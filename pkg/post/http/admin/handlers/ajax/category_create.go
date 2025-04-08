package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/app/http/response"
	"github.com/axlle-com/blog/pkg/app/logger"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *categoryController) CreateCategory(ctx *gin.Context) {
	form, formError := NewCategoryRequest().ValidateJSON(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(
				http.StatusBadRequest,
				response.Fail(http.StatusBadRequest, formError.Message, formError.Errors),
			)
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	category, err := c.categoryService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.templateProvider.GetAll()

	data := response.Body{
		"categories": categories,
		"templates":  templates,
		"category":   category,
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view":     c.RenderView("admin.category_inner", data, ctx),
				"url":      fmt.Sprintf("/admin/categories/%d", category.ID),
				"category": category,
			},
			"Запись создана",
		),
	)
}
