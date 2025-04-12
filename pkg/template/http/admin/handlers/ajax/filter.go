package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateController) FilterTemplate(ctx *gin.Context) {
	filter, validError := NeTemplateFilter().ValidateQuery(ctx)
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

	empty := &Template{}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)

	users := c.userProvider.GetAll()
	data := response.Body{
		"title":     "Страница инфо блоков",
		"template":  empty,
		"templates": templates,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"templates": templates,
				"paginator": paginator,
				"url":       filter.GetURL(),
				"view":      c.RenderView("admin.blocks_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
