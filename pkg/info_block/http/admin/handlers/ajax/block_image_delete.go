package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *blockController) DeleteBlockImage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}
	block, err := c.blockService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	err = c.blockService.DeleteImageFile(block)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}

	_, err = c.blockService.Update(block)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
}
