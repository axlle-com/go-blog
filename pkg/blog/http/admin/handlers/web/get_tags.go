package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
)

func (c *tagController) GetTags(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := models.NewTagFilter().ValidateQuery(ctx)
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

	empty := &models.PostTag{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	templates := c.template.GetAll()
	users := c.user.GetAll()

	tagsTemp, err := c.tagCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	tags := c.tagCollectionService.Aggregates(tagsTemp)

	ctx.HTML(http.StatusOK, "admin.tags", gin.H{
		"title":     "Страница тэгов",
		"tag":       empty,
		"tags":      tags,
		"templates": templates,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath()),
		},
	})
}
