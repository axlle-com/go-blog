package http

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *controller) UploadImage(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Форма не валидная!",
		})
		ctx.Abort()
		return
	}

	resource := ctx.Request.PostFormValue("resource")
	if resource == "" {
		resource = "temp"
	}

	var path string
	_, img, _ := ctx.Request.FormFile("file")
	if img != nil {
		path, err = c.uploadService.SaveUploadedFile(img, resource, user)
		if err != nil || path == "" {
			logger.Errorf("[Controller][UploadImage] Error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Изображение прикреплено",
		"data": gin.H{
			"image": path,
		},
	})
}

func (c *controller) UploadImages(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Форма не валидная!",
		})
		ctx.Abort()
		return
	}

	resource := ctx.Request.PostFormValue("resource")
	if resource == "" {
		resource = "temp"
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Errorf("[Controller][UploadImages] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	files := form.File["files"]
	paths, err := c.uploadService.SaveUploadedFiles(files, resource, user)
	if err != nil {
		logger.Error("[Controller][UploadImages] Error: paths is empty")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Изображения прикреплены",
		"data": gin.H{
			"images": paths,
		},
	})
}
