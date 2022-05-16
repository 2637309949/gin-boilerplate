package web

import (
	"gin-boilerplate/comm/db"

	"github.com/gin-gonic/gin"
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
	r.NoRoute(opt.NoRoute)
	r.Static(opt.Static.RelativePath, opt.Static.Root)

	return r
}
