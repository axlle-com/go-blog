package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPosts(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := NewPostFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL((&Post{}).AdminURL())

	templates := c.template.GetAll()
	users := c.user.GetAll()
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	postsTemp, err := c.postsService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	posts := c.postsService.GetAggregates(postsTemp)

	ctx.HTML(http.StatusOK, "admin.posts", gin.H{
		"title":        "Страница постов",
		"userProvider": user,
		"post":         &Post{},
		"posts":        posts,
		"categories":   categories,
		"templates":    templates,
		"users":        users,
		"paginator":    paginator,
		"filter":       filter,
		"menu":         models2.NewMenu(ctx.FullPath()),
	})
}
