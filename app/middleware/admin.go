package middleware

import (
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		userUUID := session.Get("user_uuid")
		userData := session.Get("user")
		if userID == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		newUser, ok := userData.(user.User)
		if !ok || !newUser.CanAdmin() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", newUser)
		c.Set("user_uuid", userUUID)
		c.Next()
	}
}
