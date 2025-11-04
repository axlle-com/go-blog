package web

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/service"
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
	api *api.Api,
) MessageWebController {
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
