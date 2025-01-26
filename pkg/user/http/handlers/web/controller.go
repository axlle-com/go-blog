package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Auth(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	Index(ctx *gin.Context)
}

func New(
	user repository.UserRepository,
) Controller {
	return &controller{
		user: user,
	}
}

type controller struct {
	*app.BaseAjax

	user repository.UserRepository
}
