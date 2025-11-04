package web

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/service"
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
	api *api.Api,
) MessageController {
	return &messageController{
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
		api:                      api,
	}
}

type messageController struct {
	*app.BaseAjax

	messageService           *service.MessageService
	messageCollectionService *service.MessageCollectionService
	api                      *api.Api
}
