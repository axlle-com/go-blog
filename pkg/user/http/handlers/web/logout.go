package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (c *controller) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{
		MaxAge: -1,
		Path:   "/",
	})
	if err := session.Save(); err != nil {
		log.Fatalln("Failed to log out")
		return
	}
	ctx.Redirect(http.StatusFound, "/")
	ctx.Abort()
	return
}
