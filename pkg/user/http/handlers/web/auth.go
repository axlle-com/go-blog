package web

import (
	. "github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) Auth(ctx *gin.Context) {
	var authInput AuthInput
	session := sessions.Default(ctx)

	if err := ctx.ShouldBind(&authInput); err != nil {
		errors := ParseBindError(err)
		for _, bindError := range errors {
			session.AddFlash(FlashErrorString(bindError))
		}
		if err := session.Save(); err != nil {
			logger.WithRequest(ctx).Error(err)
		}
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	userFound, err := c.authService.Auth(authInput)
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
			logger.WithRequest(ctx).Error(err)
		}
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}
	session.Set("user_id", userFound.ID)
	session.Set("user_uuid", userFound.UUID.String())
	session.Set("user", userFound)

	c.cache.AddUserSession(userFound.ID, session.ID())

	if err := session.Save(); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
	ctx.Abort()
	return
}
