package middles

import (
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

const (
	RequestIdKey = "X-Request-Id"
)

//Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := uuid.NewV4()
		requestId := uuid.String()
		ctx.Writer.Header().Set(RequestIdKey, requestId)
		ctx.Set(RequestIdKey, requestId)
		ctx.Next()
	}
}
