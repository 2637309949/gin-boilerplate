package web

import (
	"context"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/logger"
	"gin-boilerplate/comm/swagger/gen"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DefaultBeforeBeginFunc = func(addr string) {}
)

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

type Web struct {
	*gin.Engine
}

func (w *Web) Run(ctx context.Context, addr ...string) (err error) {
	address := resolveAddress(addr)
	srv := endless.NewServer(address, w)
	srv.BeforeBegin = DefaultBeforeBeginFunc
	go func() {
		logger.Infof(ctx, "%v/main listen and serve on 0.0.0.0%v", syscall.Getpid(), address)
		logger.Infof(ctx, "Exec `kill -1 %v` to graceful upgrade", syscall.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, err)
		}
	}()

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

	//db set up
	db.SetDsn(opt.Dialect, opt.DNS)
	db.AutoMigrate(db.GetDB())
	db.Exec(opt.Sql)

	//web init
	r := gin.New()
	r.Use(opt.Middlewares...)
	r.GET("/", opt.Index)
	r.GET(opt.Metrics, ginprom.PromHandler(promhttp.Handler()))
	r.NoRoute(opt.NoRoute)
	r.Static(opt.Static.RelativePath, opt.Static.Root)

	// gen api doc
	gn := gen.New()
	gc := gen.Config{
		SearchDir:          opt.Swagger,
		MainAPIFile:        "../main.go",
		PropNamingStrategy: "camelcase",
		MarkdownFilesDir:   "",
		OutputDir:          "./",
		ParseVendor:        true,
		ParseDependency:    true,
	}
	go gn.Build(&gc)
	return &Web{r}
}
