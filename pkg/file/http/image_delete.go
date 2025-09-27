package http

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	filePath := ctx.Param("filePath")
	if len(filePath) > 0 {
		filePath = filePath[1:]
	} else {
		logger.WithRequest(ctx).Error("[Controller][DeleteImage] Unknown path")
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	err := c.fileService.Destroy(filePath)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[Controller][Destroy] Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	err = c.uploadService.DestroyFile(filePath)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[Controller][DestroyFile] Error: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Файл успешно удален",
	})
}
