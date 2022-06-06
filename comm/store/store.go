package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"strconv"
	"time"
)

const (
	DEFAULT = time.Duration(0)
	FOREVER = time.Duration(-1)
)

var (
	PageCachePrefix = "persistence.cache"
	ErrCacheMiss    = errors.New("cache: key not found")
	ErrNotStored    = errors.New("cache: not stored")
	ErrNotSupport   = errors.New("cache: not support")
	DefaultStore    = NewInMemoryStore(time.Minute)
)

// CacheStore is the interface of a cache backend
type CacheStore interface {
	// Get retrieves an item from the cache. Returns the item or nil, and a bool indicating
	// whether the key was found.
	Get(key string, value interface{}) error

	// Set sets an item to the cache, replacing any existing item.
	Set(key string, value interface{}, expire time.Duration) error

	// Add adds an item to the cache only if an item doesn't already exist for the given
	// key, or if the existing item has expired. Returns an error otherwise.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets a new value for the cache key only if it already exists. Returns an
	// error if it does not.
	Replace(key string, data interface{}, expire time.Duration) error

	// Delete removes an item from the cache. Does nothing if the key is not in the cache.
	Delete(key string) error

	// Increment increments a real number, and returns error if the value is not real
	Increment(key string, data uint64) (uint64, error)

	// Decrement decrements a real number, and returns error if the value is not real
	Decrement(key string, data uint64) (uint64, error)

	// Flush seletes all items from the cache.
	Flush() error
}

// Serialize returns a []byte representing the passed value
func Serialize(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Deserialize deserialices the passed []byte into a the passed ptr interface{}
func Deserialize(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt
		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			}

			p.SetInt(i)
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				return err
			}

			p.SetUint(i)
			return nil
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil {
		return err
	}
	return nil
}
