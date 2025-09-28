package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *controllerCategory) GetCategories(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := models.NewCategoryFilterFilter().ValidateQuery(ctx)
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
	paginator := app.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL("/admin/post/categories")

	templates := c.templateProvider.GetAll()
	users := c.userProvider.GetAll()
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	postCategoriesTemp, err := c.categoriesService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	postCategories := c.categoriesService.GetAggregates(postCategoriesTemp)

	ctx.HTML(http.StatusOK, "admin.categories", gin.H{
		"title":          "Страница категорий",
		"postCategories": postCategories,
		"categories":     categories,
		"category":       &models.PostCategory{},
		"templates":      templates,
		"users":          users,
		"paginator":      paginator,
		"filter":         filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath()),
		},
	})
}
