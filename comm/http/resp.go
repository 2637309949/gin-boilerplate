package http

import (
	"gin-boilerplate/comm/trace"
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
	traceID, _, _ := trace.FromContext(ctx.Request.Context())
	opts = append([]Option{
		StatusOption(200),
		TraceOption(traceID),
	}, opts...)
	Responser(ctx, opts...)
}

// Fail returns Failed field
func Fail(ctx *gin.Context, opts ...Option) {
	traceID, _, _ := trace.FromContext(ctx.Request.Context())
	opts = append([]Option{
		StatusOption(500),
		TraceOption(traceID),
	}, opts...)
	Responser(ctx, opts...)
}
