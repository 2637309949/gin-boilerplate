package http

import (
	"gin-boilerplate/middles"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	Request  = http.Request
	Response = http.Response
)

func Responser(ctx *gin.Context, opts ...Option) {
	gh := gin.H{}
	for _, v := range opts {
		v(gh)
	}
	ctx.JSON(http.StatusOK, gh)
}

// Success returns successful field
func Success(ctx *gin.Context, opts ...Option) {
	requestId, _ := ctx.Get(middles.RequestIdKey)
	opts = append([]Option{
		StatusOption(200),
		TraceOption(requestId),
	}, opts...)
	Responser(ctx, opts...)
}

// Fail returns Failed field
func Fail(ctx *gin.Context, opts ...Option) {
	requestId, _ := ctx.Get(middles.RequestIdKey)
	opts = append([]Option{
		StatusOption(500),
		TraceOption(requestId),
	}, opts...)
	Responser(ctx, opts...)
}
