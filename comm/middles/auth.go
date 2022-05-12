package middles

import "github.com/gin-gonic/gin"

//AuthMiddleware ...
//JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func AuthMiddleware(auth gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth(ctx)
		ctx.Next()
	}
}
