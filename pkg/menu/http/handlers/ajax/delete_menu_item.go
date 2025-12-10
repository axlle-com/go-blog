package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
)

func (c *controllerItem) DeleteMenuItem(ctx *gin.Context) {
	id := c.GetID(ctx)

	// Если ID = 0, элемент еще не сохранен, просто возвращаем успех
	if id == 0 {
		ctx.JSON(
			http.StatusOK,
			response.OK(
				response.Body{
					"status": true,
				},
				"Элемент удален",
				nil,
			),
		)
		return
	}

	// Получаем элемент меню для получения menu_id
	menuItem, err := c.menuItemService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	menuID := menuItem.MenuID

	// Удаляем элемент меню и всех связанных потомков
	if err := c.menuItemService.Delete(id); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Получаем обновленное меню
	menu, err := c.menuService.GetByID(menuID)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Агрегируем меню для получения всех элементов
	menu, err = c.menuService.Aggregate(menu)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Рендерим обновленную страницу меню
	data := response.Body{
		"model":     menu,
		"templates": c.templates(ctx),
		"resources": app.NewResources().Resources(),
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.menu_inner", data, ctx),
			},
			"Элемент меню удален",
			nil,
		),
	)
}
