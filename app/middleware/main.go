package middleware

import (
	"github.com/axlle-com/blog/app/logger"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Main() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userUUID := session.Get("user_uuid")
		userData := session.Get("user")
		newUser, ok := userData.(user.User)
		if ok {
			ctx.Set("user", newUser)
		}

		guestUUID := session.Get("guest_uuid")
		if guestUUID == nil || guestUUID == "" {
			guestUUID = uuid.New().String()
			session.Set("guest_uuid", guestUUID)
		}

		sessionUUID := session.Get("session_uuid")
		if sessionUUID == nil || sessionUUID == "" {
			sessionUUID = uuid.New().String()
			session.Set("session_uuid", sessionUUID)
		}

		if err := session.Save(); err != nil {
			logger.Errorf("[Main][Create] Error :%v", err)
			guestUUID = ""
		}

		ctx.Set("user_uuid", userUUID)
		ctx.Set("guest_uuid", guestUUID)
		ctx.Set("session_uuid", sessionUUID)

		ctx.Next()
	}
}
