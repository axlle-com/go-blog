package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/service"
	"github.com/gin-gonic/gin"
)

type MessageController interface {
	CreateMessage(*gin.Context)
}

func NewMessageController(
	mailService *service.MailService,
) MessageController {
	return &messageController{
		mailService: mailService,
	}
}

type messageController struct {
	*app.BaseAjax

	mailService *service.MailService
}
