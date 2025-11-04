package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		return
	}

	imageId := c.getImageID(ctx)
	if imageId == 0 {
		return
	}

	_, err := c.gallery.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	image, err := c.image.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	err = c.imageService.DeleteImage(image)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, "Failed to delete file: "+err.Error(), nil),
		)
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": c.T(ctx, "ui.success.image_deleted")})
		return
	}
}
