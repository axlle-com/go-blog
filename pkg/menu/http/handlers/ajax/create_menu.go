package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/gin-gonic/gin"
)

func (c *menuController) Create(ctx *gin.Context) {
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

	menu, err := c.menuService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	menu, err = c.menuService.Aggregate(menu)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	data := response.Body{
		"model":     menu,
		"templates": c.templates(ctx),
		"resources": app.NewResources().Resources(),
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view": c.RenderView("admin.menu_inner", data, ctx),
				"url":  menu.AdminURL(),
				"menu": menu,
			},
			c.T(ctx, "ui.message.record_created"),
		),
	)
}
