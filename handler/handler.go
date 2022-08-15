package handler

import (
	"context"
	"gin-boilerplate/comm/broker"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/mark"
	"gin-boilerplate/comm/store"
	"gin-boilerplate/comm/viper"
	"net/smtp"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

//Handler...
type Handler struct {
	Store  store.CacheStore
	Broker broker.Broker
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Index
// @Description home index
// @Tags index
// @Accept  json
// @Produce  json
// @Router / [get]
func (h *Handler) Index(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "Index")()
	timemark.Mark("runtime")
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"gin_boilerplate_version": "v1.0",
		"go_version":              runtime.Version(),
	})
}

//NoRoute ...
func (h *Handler) NoRoute(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "NoRoute")()
	timemark.Mark("NoRoute")
	ctx.HTML(http.StatusOK, "notfund.tmpl", gin.H{
		"noRoute": ctx.Request.URL.Path,
	})
}

//NoRoute ...
func (h *Handler) sendEmail(ctx context.Context, el *email.Email) error {
	addr, username, identity, password, host := viper.GetString("smtp.addr"), viper.GetString("smtp.username"), viper.GetString("smtp.identity"), viper.GetString("smtp.password"), viper.GetString("smtp.host")
	err := el.Send(addr, smtp.PlainAuth(identity, username, password, host))
	return err
}
