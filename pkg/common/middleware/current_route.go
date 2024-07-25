package middleware

import "github.com/gin-gonic/gin"

func CurrentRouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRoute := c.FullPath()
		c.Set("currentRoute", currentRoute)
		c.Next()
	}
}
