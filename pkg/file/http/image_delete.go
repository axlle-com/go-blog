package http

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	filePath := ctx.Param("filePath")
	if len(filePath) > 0 {
		filePath = filePath[1:]
	} else {
		logger.Error("[Controller][DeleteImage] Не известный путь")
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	err := c.fileService.Destroy(filePath)
	if err != nil {
		logger.Errorf("[Controller][Destroy] Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	err = c.uploadService.DestroyFile(filePath)
	if err != nil {
		logger.Errorf("[Controller][DestroyFile] Error: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Файл успешно удален",
	})
}
