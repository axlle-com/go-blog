package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/gin-gonic/gin"
)

type Gallery interface {
	GetAllForResource(contracts.Resource) []contracts.Gallery
	SaveFromForm(*gin.Context) []contracts.Gallery
}

func Provider() Gallery {
	return &provider{}
}

type provider struct {
}

func (p *provider) GetAllForResource(c contracts.Resource) []contracts.Gallery {
	var collection []contracts.Gallery
	galleries, err := models.
		GalleryRepo().
		GetAllForResource(c)
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveFromForm(c *gin.Context) []contracts.Gallery {
	var collection []contracts.Gallery
	//galleries := service.SaveFromForm(c)
	//for _, gallery := range galleries {
	//	collection = append(collection, gallery)
	//}
	return collection
}
