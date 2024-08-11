package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

type Provider interface {
	GetAllForResource(contracts.Resource) []*models.Gallery
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
	logger.New().Error(err)
	return nil
}
