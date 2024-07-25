package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		MaxAge: -1,
		Path:   "/",
	})
	if err := session.Save(); err != nil {
		log.Fatalln("Failed to log out")
		return
	}
	c.Redirect(http.StatusFound, "/")
	c.Abort()
	return
}
