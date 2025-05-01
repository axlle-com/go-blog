package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		status := c.Writer.Status()

		// Решаем по заголовку Accept, является ли это API-запросом
		accept := c.GetHeader("Accept")
		isAPI := strings.Contains(accept, "application/json") ||
			strings.Contains(c.GetHeader("X-Requested-With"), "XMLHttpRequest")

		switch status {
		case http.StatusUnauthorized:
			if isAPI {
				c.JSON(status, gin.H{"message": "Не авторизован"})
			} else {
				c.HTML(status, "404", gin.H{"title": "Админка — Не авторизован"})
			}

		case http.StatusForbidden:
			if isAPI {
				c.JSON(status, gin.H{"message": "Доступ запрещён"})
			} else {
				c.HTML(status, "404", gin.H{"title": "Админка — Доступ запрещён"})
			}

		case http.StatusInternalServerError:
			errMsg := c.Errors.String()
			if isAPI {
				c.JSON(status, gin.H{"message": errMsg})
			} else {
				c.HTML(status, "404", gin.H{"error": errMsg})
			}
		}
	}
}
