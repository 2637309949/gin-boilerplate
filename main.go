package main

import (
	"context"
	"fmt"
	"gin-boilerplate/comm/broker"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/middles"
	"gin-boilerplate/comm/store"
	"gin-boilerplate/comm/viper"
	"gin-boilerplate/comm/web"
	"gin-boilerplate/handler"
	"time"
)

func main() {
	//handler
	ctx := context.TODO()
	opts := []web.OptFunc{}
	h := handler.Handler{Store: store.DefaultStore, Broker: broker.DefaultBroker}
	opts = append(opts, web.DataBase(viper.GetString("db.dialect"), viper.GetString("db.dns")))
	opts = append(opts, web.Index(h.Index))
	opts = append(opts, web.NoRoute(h.NoRoute))
	opts = append(opts, web.Static("/public", "./public"))
	opts = append(opts, web.Sql("./setup.sql"))
	opts = append(opts, web.Swagger("handler"))
	r := web.Default(opts...)

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
	r.Handle(http.MethodGet, "/api/v1/articles", middles.CachePage(h.Store, time.Minute), h.QueryArticle)
	r.Handle(http.MethodGet, "/api/v1/article/:id", middles.CachePage(h.Store, time.Minute), h.QueryArticleDetail)
	r.Handle(http.MethodPut, "/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.UpdateArticle)
	r.Handle(http.MethodDelete, "/api/v1/article/:id", middles.AuthMiddleware(h.TokenValid), h.DeleteArticle)

	//start
	r.Run(ctx, fmt.Sprintf(":%v", viper.GetString("http.port")))
}
