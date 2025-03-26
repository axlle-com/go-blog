package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/app/logger"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form, formError := NewCategoryRequest().ValidateJSON(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors":  formError.Errors,
				"message": formError.Message,
			})
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	category, err := c.categoryService.SaveFromRequest(form, found, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	categories, err := c.categoriesService.GetAllForParent(category)
	if err != nil {
		logger.Error(err)
	}

	templates := c.templateProvider.GetAll()

	data := gin.H{
		"categories": categories,
		"templates":  templates,
		"category":   category,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":     c.RenderView("admin.category_inner", data, ctx),
			"url":      fmt.Sprintf("/admin/categories/%d", category.ID),
			"category": category,
		},
	})
}
