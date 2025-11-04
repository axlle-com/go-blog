package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *postController) DeletePost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	post, err := c.postService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	if err := c.postService.PostDelete(post); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.error.server_error")})
		return
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/posts")

	users := c.api.User.GetAll()

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	postsTemp, err := c.postCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	posts := c.postCollectionService.Aggregates(postsTemp)

	data := response.Body{
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
				"view": c.RenderView("admin.posts_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
