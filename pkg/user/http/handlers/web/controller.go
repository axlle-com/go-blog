package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
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
	cache contracts.Cache,
) Controller {
	return &controller{
		userService: userService,
		authService: authService,
		cache:       cache,
	}
}

type controller struct {
	*app.BaseAjax

	cache       contracts.Cache
	userService *service.UserService
	authService *service.AuthService
}
