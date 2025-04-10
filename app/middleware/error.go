package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		switch status {
		case http.StatusUnauthorized: // 401
			c.HTML(http.StatusUnauthorized, "404", gin.H{
				"title": "Админка — Не авторизован",
			})
		case http.StatusForbidden: // 403
			c.HTML(http.StatusForbidden, "404", gin.H{
				"title": "Админка — Доступ запрещён",
			})
		case http.StatusInternalServerError: // 500
			c.HTML(http.StatusInternalServerError, "404", gin.H{
				"error": c.Errors.String(),
			})
		}
	}
}
