package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/gin-gonic/gin"
)

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	found, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	form, formError := request.NewCategoryRequest().ValidateJSON(ctx)
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

	form.PreloadFromModel(found)
	category, err := c.categoryService.SaveFromRequest(form, found, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	categories, err := c.categoriesService.GetAllForParent(category)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	infoBlocks := c.api.InfoBlock.GetAll()

	data := gin.H{
		"categories": categories,
		"templates":  c.templates(ctx),
		"category":   category,
		"collection": gin.H{
			"infoBlocks":          infoBlocks,
			"infoBlockCollection": category.InfoBlocks,
			"relationURL":         category.AdminURL(),
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":     c.RenderView("admin.category_inner", data, ctx),
			"url":      category.AdminURL(),
			"category": category,
		},
	})
}
