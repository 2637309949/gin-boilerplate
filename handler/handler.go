package handler

import (
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/http"
	"runtime"

	"github.com/gin-gonic/gin"
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

//SendEmail...
func (s *Handler) SendEmail(fromName, toAddress, toUsername, subject, textContent, token, redirctUrl, failureRedirectUrl string) error {
	return nil
}
