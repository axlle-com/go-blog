package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}

		status := ctx.Writer.Status()

		// Решаем по заголовку Accept, является ли это API-запросом
		accept := ctx.GetHeader("Accept")
		isAPI := strings.Contains(accept, "application/json") ||
			strings.Contains(ctx.GetHeader("X-Requested-With"), "XMLHttpRequest")

		switch status {
		case http.StatusUnauthorized:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Не авторизован"})
			} else {
				ctx.HTML(status, "404", gin.H{"title": "Админка — Не авторизован"})
			}

		case http.StatusForbidden:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Доступ запрещён"})
			} else {
				ctx.HTML(status, "404", gin.H{"title": "Админка — Доступ запрещён"})
			}

		case http.StatusInternalServerError:
			errMsg := ctx.Errors.String()
			if isAPI {
				ctx.JSON(status, gin.H{"message": errMsg})
			} else {
				ctx.HTML(status, "404", gin.H{"error": errMsg})
			}
		}
	}
}
