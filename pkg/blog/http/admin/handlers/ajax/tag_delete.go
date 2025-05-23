package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *tagController) Delete(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	tag, err := c.tagService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	if err := c.tagService.DeleteTags([]*PostTag{tag}); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	filter, validError := NewTagFilter().ValidateQuery(ctx)
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

	empty := &PostTag{}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	temp, err := c.tagCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	tags := c.tagCollectionService.Aggregates(temp)

	users := c.user.GetAll()
	templates := c.template.GetAll()

	data := response.Body{
		"title":     "Страница тэгов",
		"tag":       empty,
		"tags":      tags,
		"templates": templates,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.tags_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
