package ajax

import (
	"github.com/axlle-com/blog/pkg/app/http/response"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	post, err := c.post.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	if err := c.service.PostDelete(post); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

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

	users := c.user.GetAll()
	templates := c.template.GetAll()

	categories, err := c.category.GetAll()
	if err != nil {
		logger.Error(err)
	}

	posts, err := c.post.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}

	data := response.Body{
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
				"view": c.RenderView("admin.posts_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
