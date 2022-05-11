package middles

import (
	cx "gin-boilerplate/comm/util/ctx"

	"github.com/gin-gonic/gin"
)

//Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cxfr := cx.FromRequest(ctx.Request)
		rt := ctx.Request.WithContext(cxfr)
		ctx.Request = rt
		ctx.Next()
	}
}
