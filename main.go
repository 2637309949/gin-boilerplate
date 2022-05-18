package main

import (
	"context"
	"fmt"
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/gonic"
	"gin-boilerplate/comm/logger"
	"gin-boilerplate/comm/middles"
	"gin-boilerplate/comm/viper"
	"gin-boilerplate/comm/web"
	"gin-boilerplate/handler"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	//handler
	h := handler.Handler{Cache: cache.DefaultStore}
	r := web.New(web.Mode(gin.ReleaseMode),
		web.DataBase(viper.GetString("db.dialect"), viper.GetString("db.args")),
		web.Validator(new(gonic.DefaultValidator)),
		web.Middleware(ginprom.PromMiddleware(nil),
			gonic.Logger(),
			gin.Recovery(),
			middles.CORSMiddleware(),
			middles.RequestIDMiddleware(),
			gzip.Gzip(gzip.DefaultCompression)),
		web.Metrics(viper.GetString("http.metrics")),
		web.Index(h.Index),
		web.NoRoute(h.NoRoute),
		web.Static("/public", "./public"))

	//user routes
	r.POST("/api/v1/user/login", h.Login)
	r.GET("/api/v1/user/logout", h.Logout)
	r.POST("/api/v1/user/register", h.Register)
	r.POST("/api/v1/user/updatePassword", h.UpdatePassword) //for updatePassword

	r.POST("/api/v1/user/sendVerificationEmail", h.SendVerificationEmail) //for markVerified email, send verification token
	r.POST("/api/v1/user/verifyEmail", h.VerifyEmail)                     //for verify verification token and marked account

	r.POST("/api/v1/user/sendPasswordResetEmail", h.SendPasswordResetEmail) //for sendPasswordResetEmail
	r.POST("/api/v1/user/resetPassword", h.ResetPassword)                   //for resetPassword

	//auth routes
	r.POST("/api/v1/token/refresh", h.Refresh)

	//article routes
	r.POST("/api/v1/article", middles.AuthMiddleware(h.TokenValid), h.InsertArticle)
	r.GET("/api/v1/articles", middles.CachePage(h.Cache, time.Minute), h.QueryArticle)
	r.GET("/api/v1/article/:id", middles.CachePage(h.Cache, time.Minute), h.QueryArticleDetail)
	r.PUT("/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.UpdateArticle)
	r.DELETE("/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.DeleteArticle)

	//start
	logger.Infof(context.TODO(), "listen and serve on 0.0.0.0:%v", viper.GetString("http.port"))
	r.Run(fmt.Sprintf(":%v", viper.GetString("http.port")))
}
