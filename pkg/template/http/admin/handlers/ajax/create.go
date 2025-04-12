package ajax

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/pkg/template/models"
)

func (c *templateController) CreateTemplate(ctx *gin.Context) {
	form, formError := models.NewTemplateRequest().ValidateJSON(ctx)
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

	template, err := c.templateService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	data := response.Body{
		"template": template,
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view":     c.RenderView("admin.block_inner", data, ctx),
				"url":      template.AdminURL(),
				"template": template,
			},
			"Запись создана",
		),
	)
}
