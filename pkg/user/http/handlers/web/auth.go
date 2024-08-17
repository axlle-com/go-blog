package web

import (
	db "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	"github.com/axlle-com/blog/pkg/user/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
			logger.Error(err)
		}
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
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
			logger.Error(err)
		}
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	session.Set("user_id", userFound.ID)
	session.Set("user", userFound)
	sessionID := session.ID()
	cache := db.NewCache()
	cache.AddUserSession(userFound.ID, sessionID)

	if err := session.Save(); err != nil {
		logger.Error(err)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Redirect(http.StatusFound, "/admin")
	c.Abort()
	return
}
