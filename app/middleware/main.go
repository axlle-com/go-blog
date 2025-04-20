package middleware

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Main() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		guestUUID := session.Get("guest_uuid")
		if guestUUID == nil || guestUUID == "" {
			guestUUID = uuid.New().String()
			session.Set("guest_uuid", guestUUID)
		}
		if err := session.Save(); err != nil {
			logger.Errorf("[Main][Save] Error :%v", err)
			guestUUID = ""
		}

		ctx.Set("guest_uuid", guestUUID)
		ctx.Next()
	}
}
