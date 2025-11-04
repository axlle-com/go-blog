package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *categoryController) DeleteCategory(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	category, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	if err := c.categoryService.Delete(category); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	filter, validError := models.NewCategoryFilterFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(
			http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, validError.Message, validError.Errors),
		)
		ctx.Abort()
		return
	}

	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/post/categories")

	users := c.api.User.GetAll()

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	postCategoriesTemp, err := c.categoriesService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	postCategories := c.categoriesService.GetAggregates(postCategoriesTemp)

	data := response.Body{
		"title":          "Страница постов",
		"category":       &models.PostCategory{},
		"categories":     categories,
		"postCategories": postCategories,
		"templates":      c.templates(ctx),
		"users":          users,
		"paginator":      paginator,
		"filter":         filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.categories_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
