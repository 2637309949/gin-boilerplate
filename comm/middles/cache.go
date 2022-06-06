package middles

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"gin-boilerplate/comm/logger"
	mStore "gin-boilerplate/comm/store"

	"github.com/gin-gonic/gin"
)

var PageCachePrefix = "uy7hxnw8"

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

type cachedWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	store   mStore.CacheStore
	expire  time.Duration
	key     string
}

// RegisterResponseCacheGob registers the responseCache type with the encoding/gob package
func RegisterResponseCacheGob() {
	gob.Register(responseCache{})
}

// CreateKey creates a package specific key for a given string
func CreateKey(u string) string {
	return urlEscape(PageCachePrefix, u)
}

func urlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}

func newCachedWriter(store mStore.CacheStore, expire time.Duration, writer gin.ResponseWriter, key string) *cachedWriter {
	return &cachedWriter{writer, 0, false, store, expire, key}
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *cachedWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		store := w.store
		var cache responseCache
		if err := store.Get(w.key, &cache); err == nil {
			data = append(cache.Data, data...)
		}
		//cache responses with a status code < 300
		if w.Status() < 300 {
			val := responseCache{
				w.Status(),
				w.Header(),
				data,
			}
			err = store.Set(w.key, val, w.expire)
			if err != nil {
				log.Println(context.TODO(), err)
			}
		}
	}
	return ret, err
}

func (w *cachedWriter) WriteString(data string) (n int, err error) {
	ret, err := w.ResponseWriter.WriteString(data)
	//cache responses with a status code < 300
	if err == nil && w.Status() < 300 {
		store := w.store
		val := responseCache{
			w.Status(),
			w.Header(),
			[]byte(data),
		}
		store.Set(w.key, val, w.expire)
	}
	return ret, err
}

// CachePage Decorator
func CachePage(store mStore.CacheStore, expire time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rsp responseCache
		url := c.Request.URL
		key := CreateKey(url.RequestURI())
		if err := store.Get(key, &rsp); err != nil {
			if err != mStore.ErrCacheMiss {
				logger.Error(c.Request.Context(), err.Error())
			}
			writer := newCachedWriter(store, expire, c.Writer, key)
			c.Writer = writer
			c.Next()

			if c.IsAborted() {
				store.Delete(key)
			}
		} else {
			logger.Warnf(c.Request.Context(), "Load from cache")
			c.Writer.WriteHeader(rsp.Status)
			for k, vals := range rsp.Header {
				for _, v := range vals {
					c.Writer.Header().Set(k, v)
				}
			}
			c.Writer.Write(rsp.Data)
			c.Abort()
		}
	}
}

// Cache middles
func Cache(store mStore.CacheStore, time time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		CachePage(store, time)(ctx)
	}
}
