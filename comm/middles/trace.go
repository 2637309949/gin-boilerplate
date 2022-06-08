package middles

import (
	cx "gin-boilerplate/comm/util/ctx"

	"github.com/gin-gonic/gin"
)

//Generate a unique ID and attach it to each request for future reference or use
func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(cx.FromRequest(ctx.Request))
		ctx.Next()
	}
}
