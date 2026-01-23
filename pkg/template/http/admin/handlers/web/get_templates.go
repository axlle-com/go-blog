package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/http/request"
	template "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *templateWebController) GetTemplates(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := request.NeTemplateFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.error.server_error")})
		return
	}

	empty := &template.Template{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	users := c.api.User.GetAll()

	temp, err := c.templateCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	templates := c.templateCollectionService.Aggregates(temp)
	resources := app.NewResources()
	c.RenderHTML(ctx, http.StatusOK, "admin.templates", gin.H{
		"title":         c.T(ctx, "ui.page.templates"),
		"templateModel": empty,
		"resources":     resources.Resources(),
		"themes":        resources.Themes(),
		"templates":     templates,
		"users":         users,
		"paginator":     paginator,
		"filter":        filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
		},
	})
}
