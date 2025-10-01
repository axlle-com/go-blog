package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *tagController) Filter(ctx *gin.Context) {
	filter, validError := models.NewTagFilter().ValidateQuery(ctx)
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

	empty := &models.PostTag{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	temp, err := c.tagCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	tags := c.tagCollectionService.Aggregates(temp)

	templates := c.template.GetAll()
	users := c.user.GetAll()

	data := gin.H{
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
				"tags":      tags,
				"paginator": paginator,
				"url":       filter.GetURL(),
				"view":      c.RenderView("admin.tags_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
