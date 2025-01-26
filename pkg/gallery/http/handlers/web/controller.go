package web

import (
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller interface {
	DeleteImage(*gin.Context)
}

func New(
	gallery repository.GalleryRepository,
	image repository.GalleryImageRepository,
	imageService *service.ImageService,
) Controller {
	return &controller{
		gallery:      gallery,
		image:        image,
		imageService: imageService,
	}
}

type controller struct {
	gallery      repository.GalleryRepository
	image        repository.GalleryImageRepository
	imageService *service.ImageService
}

func (c *controller) GetID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
	}
	return uint(id)
}

func (c *controller) getImageID(ctx *gin.Context) uint {
	idParam := ctx.Param("image_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
	}
	return uint(id)
}
