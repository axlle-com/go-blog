package ajax

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *categoryController) DeleteCategoryImage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}
	category, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	err = c.categoryService.DeleteImageFile(category)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}

	newCategory := *category
	_, err = c.categoryService.Update(&newCategory, category, c.GetUser(ctx))
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
	return
}
