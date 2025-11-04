package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *menuController) GetMenus(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := request.NewMenuFilter().ValidateQuery(ctx)
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
	empty := &models.Menu{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	menus, err := c.menuCollectionService.WithPaginate(paginator, filter.ToFilter())
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	c.RenderHTML(ctx, http.StatusOK, "admin.menus", gin.H{
		"title":     "Страница меню",
		"menus":     menus,
		"model":     empty,
		"templates": c.templates(ctx),
		"paginator": paginator,
		"filter":    filter,
		"settings": gin.H{
			"title":     "Страница меню",
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      models.NewMenu(ctx.FullPath(), c.BuildT(ctx)),
		},
	})
}
