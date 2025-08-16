package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *controller) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	guestUUID := session.Get("guest_uuid")
	if guestUUID == nil || guestUUID == "" {
		guestUUID = uuid.New().String()
	}

	session.Clear()

	session.Set("guest_uuid", guestUUID)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7 дней
		HttpOnly: true,
	})
	if err := session.Save(); err != nil {
		logger.WithRequest(ctx).Error("Failed to log out")
	}
	ctx.Redirect(http.StatusFound, "/")
	ctx.Abort()
}
