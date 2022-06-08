package middles

import "github.com/gin-gonic/gin"

//AuthMiddleware ...
func AuthMiddleware(auth gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth(ctx)
	}
}
