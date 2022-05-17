package middles

import "github.com/gin-gonic/gin"

//AuthMiddleware ...
func AuthMiddleware(authMiddle gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authMiddle(ctx)
	}
}
