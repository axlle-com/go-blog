package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *postController) GetPosts(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := models.NewPostFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.message.server_error")})
		return
	}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL((&models.Post{}).AdminURL())

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

	c.RenderHTML(ctx, http.StatusOK, "admin.posts", gin.H{
		"title":      c.T(ctx, "ui.name.posts"),
		"post":       &models.Post{},
		"posts":      posts,
		"categories": categories,
		"templates":  c.templates(ctx),
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
		},
	})
}
