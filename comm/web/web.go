package web

import (
	"context"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/gonic"
	"gin-boilerplate/comm/logger"
	"gin-boilerplate/comm/middles"
	"gin-boilerplate/comm/viper"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DefaultBeforeBeginFunc = func(addr string) {}
)

func resolveAddress(addr ...string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		if strings.Contains(addr[0], ":") {
			return addr[0]
		}
		return ":" + addr[0]
	default:
		panic("too many parameters")
	}
}

type Web struct {
	*gin.Engine
}

func (w *Web) listen(ctx context.Context, srv *http.Server) (err error) {
	https := viper.GetBool("listen.https")
	if https {
		certFile, keyFile := viper.GetString("listen.certFile"), viper.GetString("listen.keyFile")
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, err)
			return err
		}
	} else {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, err)
			return err
		}
	}
	return nil
}

func (w *Web) Run(ctx context.Context) (err error) {
	address := resolveAddress(viper.GetString("listen.port"))
	srv := NewServer(address, w)
	logger.Infof(ctx, "%v listen and serve %v/main on 0.0.0.0%v", viper.GetString("name"), syscall.Getpid(), address)
	logger.Infof(ctx, "exec `kill -1 %v` to graceful upgrade", syscall.Getpid())
	go w.listen(ctx, &srv.Server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info(ctx, "Gracefully shutdown Server ...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, err)
		return err
	}
	logger.Info(ctx, "Server exiting")

	return nil
}
func New(opts ...OptFunc) *Web {
	var opt Option
	for _, v := range opts {
		v(&opt)
	}

	// Dialect set up
	if len(opt.Dialect) > 0 {
		db.SetDsn(opt.Dialect, opt.DNS)
		db.AutoMigrate(db.GetDB())
	}
	if len(opt.Sql) > 0 {
		db.Exec(opt.Sql)
	}

	// web init
	r := gin.New()
	r.Use(opt.Middlewares...)
	r.GET("/", opt.Index)
	r.GET(opt.Metrics, ginprom.PromHandler(promhttp.Handler()))
	r.NoRoute(opt.NoRoute)
	if len(opt.Static.RelativePath) > 0 {
		r.Static(opt.Static.RelativePath, opt.Static.Root)
	}
	if len(opt.Template) > 0 {
		r.LoadHTMLGlob(opt.Template)
	}

	//Build builds swagger json file  for given searchDir and mainAPIFile. Returns json
	GenerageSwaggerDoc(opt.Swagger)
	return &Web{r}
}

func Default(opts ...OptFunc) *Web {
	opts = append(opts, Mode(gin.ReleaseMode))
	opts = append(opts, Validator(new(gonic.DefaultValidator)))
	opts = append(opts, Metrics("/metrics"))
	opts = append(opts, Middleware(gzip.Gzip(gzip.DefaultCompression), middles.TraceMiddleware(), gin.Recovery(), gonic.Logger(), middles.CORSMiddleware(), ginprom.PromMiddleware(nil)))
	return New(opts...)
}
