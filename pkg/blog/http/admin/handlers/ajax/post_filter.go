package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *postController) FilterPosts(ctx *gin.Context) {
	filter, validError := models.NewPostFilter().ValidateQuery(ctx)
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

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/posts")

	postsTemp, err := c.postCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	posts := c.postCollectionService.Aggregates(postsTemp)

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	users := c.api.User.GetAll()
	data := gin.H{
		"title":      "Страница постов",
		"post":       &models.Post{},
		"posts":      posts,
		"categories": categories,
		"templates":  c.templates(ctx),
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"posts":     posts,
				"paginator": paginator,
				"url":       filter.GetURL(),
				"view":      c.RenderView("admin.posts_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
