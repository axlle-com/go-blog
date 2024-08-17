package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/gin-gonic/gin"
)

type Provider interface {
	GetAllForResource(contracts.Resource) []*models.Gallery
	SaveFromForm(*gin.Context) []*models.Gallery
}

func NewProvider() Provider {
	return &providerGallery{}
}

type providerGallery struct {
}

func (p *providerGallery) GetAllForResource(c contracts.Resource) []*models.Gallery {
	galleries, err := models.
		NewGalleryRepository().
		GetAllForResource(c)
	if err == nil {
		return galleries
	}
	logger.Error(err)
	return nil
}

func (p *providerGallery) SaveFromForm(c *gin.Context) []*models.Gallery {
	return service.SaveFromForm(c)
}
