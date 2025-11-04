package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
)

func (c *templateController) FilterTemplate(ctx *gin.Context) {
	filter, validError := request.NeTemplateFilter().ValidateQuery(ctx)
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

	empty := &models.Template{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)

	users := c.api.User.GetAll()
	resources := app.NewResources()
	data := response.Body{
		"title":         "Страница шаблонов",
		"templateModel": empty,
		"templates":     templates,
		"users":         users,
		"paginator":     paginator,
		"filter":        filter,
		"resources":     resources.Resources(),
		"themes":        resources.Themes(),
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"templates": templates,
				"paginator": paginator,
				"url":       filter.GetURL(),
				"view":      c.RenderView("admin.templates_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
