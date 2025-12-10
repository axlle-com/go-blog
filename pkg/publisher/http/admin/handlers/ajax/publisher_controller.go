package ajax

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/publisher/service"
	"github.com/gin-gonic/gin"
)

type PublisherController interface {
	Filter(*gin.Context)
}

func NewPublisherController(
	collectionService *service.CollectionService,
	api *api.Api,
) PublisherController {
	return &controller{
		collectionService: collectionService,
		api:               api,
	}
}

type controller struct {
	*app.BaseAjax

	collectionService *service.CollectionService
	api               *api.Api
}
