package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/gin-gonic/gin"
)

func (c *controllerMenu) Update(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	found, err := c.menuService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	form, formError := request.NewMenuRequest().ValidateJSON(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(
				http.StatusBadRequest,
				response.Fail(http.StatusBadRequest, formError.Message, formError.Errors),
			)
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	menu, err := c.menuService.SaveFromRequest(form, found, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	menu, err = c.menuService.Aggregate(menu)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "admin.404", gin.H{"title": err.Error()})
		return
	}

	templates, err := c.templateProvider.GetForResources(menu)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	publishers, err := c.postProvider.GetPublishers()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	data := response.Body{
		"model":      menu,
		"templates":  templates,
		"publishers": publishers,
		"resources":  app.NewResources().Resources(),
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view": c.RenderView("admin.menu_inner", data, ctx),
				"url":  menu.AdminURL(),
				"menu": menu,
			},
			"Запись создана",
		),
	)
}
