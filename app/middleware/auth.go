package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("user_id")
		userUUID := session.Get("user_uuid")
		user := session.Get("user")
		if userID == nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Set("user_uuid", userUUID)
		ctx.Next()
	}
}
