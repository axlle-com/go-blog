package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/gin-gonic/gin"
)

func (c *categoryController) CreateCategory(ctx *gin.Context) {
	form, formError := request.NewCategoryRequest().ValidateJSON(ctx)
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

	categories, err := c.categoriesService.GetAllForParent(category)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	infoBlocks := c.api.InfoBlock.GetAll()

	data := response.Body{
		"categories": categories,
		"templates":  c.templates(ctx),
		"category":   category,
		"collection": gin.H{
			"infoBlocks":          infoBlocks,
			"infoBlockCollection": category.InfoBlocks,
			"relationURL":         category.AdminURL(),
		},
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view":     c.RenderView("admin.category_inner", data, ctx),
				"url":      category.AdminURL(),
				"category": category,
			},
			c.T(ctx, "ui.message.record_created"),
		),
	)
}
