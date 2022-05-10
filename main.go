package main

import (
	"time"

	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/handler"
	"gin-boilerplate/middles"
	"gin-boilerplate/types"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	//setup db
	db.SetDsn("sqlite3", "./sqlite.db")
	db.AutoMigrate(db.GetDB())
	gin.SetMode(gin.ReleaseMode)
	binding.Validator = new(types.DefaultValidator)
	hdl := handler.Handler{Cache: cache.DefaultStore}

	//http handler
	r := gin.New()
	r.Use(gin.Logger())
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
	r.GET("/api/v1/articles", middles.CachePage(cache.DefaultStore, time.Minute), hdl.QueryArticle)
	r.GET("/api/v1/article/:id", middles.TokenAuthMiddleware(hdl.TokenValid), middles.CachePage(cache.DefaultStore, time.Minute), hdl.QueryArticleDetail)
	r.PUT("/api/v1/article/:id", middles.TokenAuthMiddleware(hdl.TokenValid), hdl.UpdateArticle)
	r.DELETE("/api/v1/article/:id", middles.TokenAuthMiddleware(hdl.TokenValid), hdl.DeleteArticle)

	r.Run(":8080")
}
