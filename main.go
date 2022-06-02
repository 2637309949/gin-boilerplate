package main

import (
	"context"
	"fmt"
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/gonic"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/middles"
	"gin-boilerplate/comm/viper"
	"gin-boilerplate/comm/web"
	"gin-boilerplate/handler"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	//handler
	todo := context.TODO()
	h := handler.Handler{Cache: cache.DefaultStore}
	r := web.New(web.Mode(gin.ReleaseMode),
		web.DataBase(viper.GetString("db.dialect"), viper.GetString("db.dns")),
		web.Validator(new(gonic.DefaultValidator)),
		web.Middleware(ginprom.PromMiddleware(nil),
			gonic.Logger(),
			gin.Recovery(),
			middles.CORSMiddleware(),
			middles.RequestIDMiddleware(),
			gzip.Gzip(gzip.DefaultCompression)),
		web.Metrics("/metrics"),
		web.Index(h.Index),
		web.NoRoute(h.NoRoute),
		web.Static("/public", "./public"),
		web.Sql("./setup.sql"),
		web.Swagger("handler"))

	//User routes
	r.Handle(http.MethodPost, "/api/v1/user/login", h.Login)
	r.Handle(http.MethodGet, "/api/v1/user/logout", h.Logout)
	r.Handle(http.MethodPost, "/api/v1/user/register", h.Register)
	r.Handle(http.MethodPost, "/api/v1/user/updatePassword", h.UpdatePassword) //for updatePassword
	//
	r.Handle(http.MethodPost, "/api/v1/user/sendVerificationEmail", h.SendVerificationEmail) //for markVerified email, send verification token
	r.Handle(http.MethodPost, "/api/v1/user/verifyEmail", h.VerifyEmail)                     //for verify verification token and marked account
	//
	r.Handle(http.MethodPost, "/api/v1/user/sendPasswordResetEmail", h.SendPasswordResetEmail) //for sendPasswordResetEmail
	r.Handle(http.MethodPost, "/api/v1/user/resetPassword", h.ResetPassword)                   //for resetPassword

	//Auth routes
	r.Handle(http.MethodPost, "/api/v1/token/refresh", h.Refresh)

	//Article routes
	r.Handle(http.MethodPost, "/api/v1/article", middles.AuthMiddleware(h.TokenValid), h.InsertArticle)
	r.Handle(http.MethodGet, "/api/v1/articles", middles.CachePage(h.Cache, time.Minute), h.QueryArticle)
	r.Handle(http.MethodGet, "/api/v1/article/:id", middles.CachePage(h.Cache, time.Minute), h.QueryArticleDetail)
	r.Handle(http.MethodPut, "/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.UpdateArticle)
	r.Handle(http.MethodDelete, "/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.DeleteArticle)

	//start
	r.Run(todo, fmt.Sprintf(":%v", viper.GetString("http.port")))
}
