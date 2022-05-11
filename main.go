package main

import (
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/gonic"
	"gin-boilerplate/comm/middles"
	"gin-boilerplate/handler"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	//setup db
	db.Init()
	gin.SetMode(gin.ReleaseMode)
	binding.Validator = new(gonic.DefaultValidator)
	hdl := handler.Handler{Cache: cache.DefaultStore}

	//http handler
	r := gin.New()
	r.Use(gonic.Logger())
	r.Use(gin.Recovery())
	r.Use(middles.CORSMiddleware())
	r.Use(middles.RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Static("/public", "./public")

	//default routes
	r.GET("/", hdl.Index)
	r.NoRoute(hdl.NoRoute)

	//users routes
	r.POST("/api/v1/user/login", hdl.Login)
	r.POST("/api/v1/user/register", hdl.Register)
	r.GET("/api/v1/user/logout", hdl.Logout)

	//auth routes
	r.POST("/api/v1/token/refresh", hdl.Refresh)

	//article routes
	r.POST("/api/v1/article", middles.TokenAuthMiddleware(hdl.TokenValid), hdl.InsertArticle)
	r.GET("/api/v1/articles", middles.CachePage(hdl.Cache, time.Minute), hdl.QueryArticle)
	r.GET("/api/v1/article/:id", middles.CachePage(hdl.Cache, time.Minute), hdl.QueryArticleDetail)
	r.PUT("/api/v1/article/:id", middles.TokenAuthMiddleware(hdl.TokenValid), hdl.UpdateArticle)
	r.DELETE("/api/v1/article/:id", middles.TokenAuthMiddleware(hdl.TokenValid), hdl.DeleteArticle)

	//start
	r.Run(":8080")
}
