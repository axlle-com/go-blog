package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/service"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type MessageController interface {
	GetMessage(*gin.Context)
	GetMessages(*gin.Context)
	DeleteMessage(ctx *gin.Context)
	CreateMessage(*gin.Context)
}

func NewMessageController(
	messageService *service.MessageService,
	messageCollectionService *service.MessageCollectionService,
	userProvider userProvider.UserProvider,
) MessageController {
	return &messageController{
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
		userProvider:             userProvider,
	}
}

type messageController struct {
	*app.BaseAjax

	messageService           *service.MessageService
	messageCollectionService *service.MessageCollectionService
	userProvider             userProvider.UserProvider
}
