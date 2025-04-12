package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateWebController) GetTemplates(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := NeTemplateFilter().ValidateQuery(ctx)
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

	empty := &Template{}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	users := c.userProvider.GetAll()

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)

	ctx.HTML(http.StatusOK, "admin.blocks", gin.H{
		"title":     "Страница инфо блоков",
		"user":      user,
		"template":  empty,
		"templates": templates,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
		"menu":      models2.NewMenu(ctx.FullPath()),
	})
}
