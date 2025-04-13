package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateController) DeleteTemplate(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	block, err := c.templateService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	if err := c.templateService.Delete(block); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}

	empty := &Template{}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	users := c.userProvider.GetAll()

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)

	data := response.Body{
		"title":         "Страница шаблонов",
		"templateModel": empty,
		"templates":     templates,
		"users":         users,
		"paginator":     paginator,
		"filter":        filter,
		"resources":     models.NewResource().Resources(),
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.templates_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
