package middles

import "github.com/gin-gonic/gin"

//TokenAuthMiddleware ...
//JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware(auth gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth(ctx)
		ctx.Next()
	}
}
