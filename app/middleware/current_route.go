package middleware

import "github.com/gin-gonic/gin"

func CurrentRouteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentRoute := ctx.FullPath()
		ctx.Set("currentRoute", currentRoute)
		ctx.Next()
	}
}
