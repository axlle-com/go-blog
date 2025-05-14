package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *postController) FilterPosts(ctx *gin.Context) {
	filter, validError := NewPostFilter().ValidateQuery(ctx)
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

	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
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

	templates := c.template.GetAll()
	users := c.user.GetAll()
	data := gin.H{
		"title":      "Страница постов",
		"post":       &Post{},
		"posts":      posts,
		"categories": categories,
		"templates":  templates,
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
