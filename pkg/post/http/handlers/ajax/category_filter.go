package ajax

import (
	"github.com/axlle-com/blog/pkg/app/http/response"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *categoryController) FilterCategory(ctx *gin.Context) {
	filter, validError := NewCategoryFilterFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(
			http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, validError.Message, validError.Errors),
		)
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Message("Ошибка сервера"))
		return
	}

	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	postCategories, err := c.categoriesService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.templateProvider.GetAll()
	users := c.userProvider.GetAll()
	data := response.Body{
		"title":          "Страница категорий",
		"postCategories": postCategories,
		"categories":     categories,
		"templates":      templates,
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
