package web

import (
	. "github.com/axlle-com/blog/pkg/common/errors"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	"github.com/axlle-com/blog/pkg/user/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Auth(c *gin.Context) {
	var authInput AuthInput
	session := sessions.Default(c)

	if err := c.ShouldBind(&authInput); err != nil {
		errors := ParseBindError(err)
		for _, bindError := range errors {
			session.AddFlash(FlashErrorString(bindError))
		}
		if err := session.Save(); err != nil {
			log.Println(err)
		}
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userFound, err := service.Auth(authInput)
	if err != nil {
		session.AddFlash(
			FlashErrorString(
				BindError{
					Field:   GeneralFieldName,
					Message: err.Error(),
				},
			),
		)
		err := session.Save()
		if err != nil {
			log.Println(err)
		}
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Set("user_id", userFound.ID)
	if err := session.Save(); err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Redirect(http.StatusFound, "/admin")
}
