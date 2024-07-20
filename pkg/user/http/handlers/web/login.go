package web

import (
	. "github.com/axlle-com/blog/pkg/common/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	flashes := session.Flashes()
	errorMessages := ParseFlashes(flashes)
	err := session.Save()
	if err != nil {
		log.Println(err)
	}
	c.HTML(
		http.StatusOK,
		"admin.login",
		gin.H{
			"Title":  "Авторизация",
			"Errors": errorMessages,
		},
	)
}
