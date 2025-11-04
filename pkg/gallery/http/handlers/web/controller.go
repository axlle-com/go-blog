package web

import (
	"net/http"
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/gin-gonic/gin"
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
	*app.BaseAjax

	gallery      repository.GalleryRepository
	image        repository.GalleryImageRepository
	imageService *service.ImageService
}

func (c *controller) GetID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
	}
	return uint(id)
}

func (c *controller) getImageID(ctx *gin.Context) uint {
	idParam := ctx.Param("image_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
	}
	return uint(id)
}
