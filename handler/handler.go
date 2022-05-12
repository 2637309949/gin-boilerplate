package handler

import (
	"context"
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/http"
	"net/smtp"
	"runtime"

	"gin-boilerplate/comm/viper"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

type Handler struct {
	Cache cache.CacheStore
}

//Index ...
func (s *Handler) Index(ctx *gin.Context) {
	tk := gin.H{
		"ginBoilerplateVersion": "v1.0",
		"goVersion":             runtime.Version(),
	}
	http.Success(ctx, http.FlatOption(tk))
}

//NoRoute ...
func (s *Handler) NoRoute(ctx *gin.Context) {
	http.Fail(ctx, http.MsgOption("NotFound"), http.StatusOption(http.StatusNotFound))
}

//NoRoute ...
func (s *Handler) sendEmail(ctx context.Context, el *email.Email) error {
	addr, username, identity, password, host := viper.GetString("smtp.addr"), viper.GetString("smtp.username"), viper.GetString("smtp.identity"), viper.GetString("smtp.password"), viper.GetString("smtp.host")
	err := el.Send(addr, smtp.PlainAuth(identity, username, password, host))
	return err
}
