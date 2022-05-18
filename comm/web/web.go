package web

import (
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/swagger/gen"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(opts ...OptFunc) *gin.Engine {
	var opt Option
	for _, v := range opts {
		v(&opt)
	}

	//db set up
	db.SetDsn(opt.Dialect, opt.DialectArgs...)
	db.AutoMigrate(db.GetDB())

	//web init
	r := gin.New()
	r.Use(opt.Middlewares...)
	r.GET("/", opt.Index)
	r.GET(opt.Metrics, ginprom.PromHandler(promhttp.Handler()))
	r.NoRoute(opt.NoRoute)
	r.Static(opt.Static.RelativePath, opt.Static.Root)

	// gen api doc
	gn := gen.New()
	genCfg := gen.Config{
		SearchDir:          opt.Swagger,
		MainAPIFile:        "main.go",
		PropNamingStrategy: "camelcase",
		MarkdownFilesDir:   "",
		OutputDir:          "./",
		ParseVendor:        true,
		ParseDependency:    true,
	}
	go gn.Build(&genCfg)
	return r
}
