package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/user/http/models"
)

func (c *controller) Auth(ctx *gin.Context) {
	var authInput models.AuthInput
	session := sessions.Default(ctx)

	if err := ctx.ShouldBind(&authInput); err != nil {
		errors := errutil.NewBindError(err)
		for _, bindError := range errors {
			session.AddFlash(errutil.FlashErrorString(bindError))
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
			errutil.FlashErrorString(
				errutil.BindError{
					Field:   errutil.GeneralFieldName,
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
}
