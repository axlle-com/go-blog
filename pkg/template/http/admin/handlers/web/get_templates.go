package web

import (
	"github.com/axlle-com/blog/app/logger"
	mApp "github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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
	paginator := mApp.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	users := c.userProvider.GetAll()

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)

	ctx.HTML(http.StatusOK, "admin.templates", gin.H{
		"title":         "Страница шаблонов",
		"templateModel": empty,
		"resources":     mApp.NewResources().Resources(),
		"templates":     templates,
		"users":         users,
		"paginator":     paginator,
		"filter":        filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      models2.NewMenu(ctx.FullPath()),
		},
	})
}
