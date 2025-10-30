package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/gin-gonic/gin"
)

func (c *templateController) CreateTemplate(ctx *gin.Context) {
	form, formError := request.NewTemplateRequest().ValidateJSON(ctx)
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

	resources := app.NewResources()
	data := response.Body{
		"templateModel": template,
		"resources":     resources.Resources(),
		"themes":        resources.Themes(),
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view":     c.RenderView("admin.template_inner", data, ctx),
				"url":      template.AdminURL(),
				"template": template,
			},
			"Запись создана",
		),
	)
}
