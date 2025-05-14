package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	. "github.com/axlle-com/blog/pkg/blog/http/admin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *tagController) Create(ctx *gin.Context) {
	form, formError := NewTagRequest().ValidateJSON(ctx)
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

	tag, err := c.tagService.SaveFromRequest(form, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	templates := c.template.GetAll()
	infoBlocks := c.infoBlock.GetAll()

	data := response.Body{
		"templates": templates,
		"tag":       tag,
		"collection": gin.H{
			"infoBlocks":         infoBlocks,
			"ifoBlockCollection": tag.InfoBlocks,
			"relationURL":        tag.AdminURL(),
		},
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view": c.RenderView("admin.tag_inner", data, ctx),
				"url":  tag.AdminURL(),
				"tag":  tag,
			},
			"Запись создана",
		),
	)
}
