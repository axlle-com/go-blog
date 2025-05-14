package ajax

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *tagController) DeleteImage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}
	post, err := c.tagService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	err = c.tagService.DeleteImageFile(post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}

	_, err = c.tagService.Update(post)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
	return
}
