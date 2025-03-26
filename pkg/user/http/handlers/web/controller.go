package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/user/service"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Auth(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	Index(ctx *gin.Context)
}

func NewUserWebController(
	userService *service.UserService,
	authService *service.AuthService,
) Controller {
	return &controller{
		userService: userService,
		authService: authService,
	}
}

type controller struct {
	*app.BaseAjax

	userService *service.UserService
	authService *service.AuthService
}
