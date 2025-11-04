package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *categoryController) FilterCategory(ctx *gin.Context) {
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Message(c.T(ctx, "ui.error.server_error")))
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
		"title":          "Страница категорий",
		"category":       &models.PostCategory{},
		"postCategories": postCategories,
		"categories":     categories,
		"templates":      c.templates(ctx),
		"users":          users,
		"paginator":      paginator,
		"filter":         filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"postCategories": postCategories,
				"paginator":      paginator,
				"url":            filter.GetURL(),
				"view":           c.RenderView("admin.categories_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
