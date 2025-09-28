package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/axlle-com/blog/app/logger"
	user "github.com/axlle-com/blog/pkg/user/models"
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
		if err := session.Save(); err != nil {
			logger.Errorf("[Main][Save] Error :%v", err)
			guestUUID = ""
		}

		ctx.Set("user_uuid", userUUID)
		ctx.Set("guest_uuid", guestUUID)
		ctx.Next()
	}
}
