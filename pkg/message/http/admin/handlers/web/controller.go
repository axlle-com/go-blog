package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/service"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type MessageWebController interface {
	GetMessage(*gin.Context)
	GetMessages(*gin.Context)
	CreateMessage(*gin.Context)
}

func NewMessageWebController(
	messageService *service.MessageService,
	messageCollectionService *service.MessageCollectionService,
	userProvider userProvider.UserProvider,
) MessageWebController {
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
