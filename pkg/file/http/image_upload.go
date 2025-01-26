package http

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) UploadImage(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Форма не валидная!",
		})
		ctx.Abort()
		return
	}

	var path string
	_, img, _ := ctx.Request.FormFile("file")
	if img != nil {
		path, err = c.service.SaveUploadedFile(img, "temp")
		if err != nil || path == "" {
			logger.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Изображение загружено",
		"data": gin.H{
			"image": path,
		},
	})
}

func (c *controller) UploadImages(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Форма не валидная!",
		})
		ctx.Abort()
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	files := form.File["files"]
	paths := c.service.SaveUploadedFiles(files, "temp")
	if len(paths) <= 0 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Изображения загружены",
		"data": gin.H{
			"images": paths,
		},
	})
}
