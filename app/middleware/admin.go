package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/pkg/user/models"
)

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("user_id")
		userUUID := session.Get("user_uuid")
		userData := session.Get("user")
		if userID == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		newUser, ok := userData.(models.User)
		if !ok || !newUser.CanAdmin() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", newUser)
		ctx.Set("user_uuid", userUUID)
		ctx.Next()
	}
}
