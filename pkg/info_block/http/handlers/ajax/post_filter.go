package ajax

import (
	"github.com/axlle-com/blog/pkg/app/http/response"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) FilterPosts(ctx *gin.Context) {
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

	paginator := models.NewPaginator(ctx.Request.URL.Query())
	paginator.AddQueryString(string(filter.GetQueryString()))
	posts, err := c.post.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	categories, err := c.category.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
	users := c.user.GetAll()
	data := gin.H{
		"title":      "Страница постов",
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
