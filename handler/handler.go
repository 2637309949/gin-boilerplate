package handler

import (
	"context"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/store"
	"gin-boilerplate/comm/viper"
	"net/smtp"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

type Handler struct {
	Store store.CacheStore
}

//Index ...
func (h *Handler) Index(ctx *gin.Context) {
	tk := gin.H{
		"gin_boilerplate_version": "v1.0",
		"go_version":              runtime.Version(),
	}
	http.Success(ctx, http.FlatOption(tk))
}

//NoRoute ...
func (h *Handler) NoRoute(ctx *gin.Context) {
	http.Fail(ctx, http.MsgOption("NotFound"), http.StatusOption(http.StatusNotFound))
}

//NoRoute ...
func (h *Handler) sendEmail(ctx context.Context, el *email.Email) error {
	addr, username, identity, password, host := viper.GetString("smtp.addr"), viper.GetString("smtp.username"), viper.GetString("smtp.identity"), viper.GetString("smtp.password"), viper.GetString("smtp.host")
	err := el.Send(addr, smtp.PlainAuth(identity, username, password, host))
	return err
}
