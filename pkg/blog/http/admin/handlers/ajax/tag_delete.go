package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *tagController) Delete(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	tag, err := c.tagService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	if err := c.tagService.DeleteTags([]*models.PostTag{tag}); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.error.server_error")})
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

	users := c.api.User.GetAll()

	data := response.Body{
		"title":     c.T(ctx, "ui.page.tags"),
		"tag":       empty,
		"tags":      tags,
		"templates": c.templates(ctx),
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
			c.T(ctx, "ui.success.record_deleted"),
			paginator,
		),
	)
}
