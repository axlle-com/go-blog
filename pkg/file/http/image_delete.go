package http

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	filePath := ctx.Param("filePath")
	if len(filePath) > 0 {
		filePath = filePath[1:]
	} else {
		logger.Error(errors.New("Не известный путь"))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// TODO
	ctx.JSON(200, gin.H{
		"message": "Файл успешно удален",
	})
}
