package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/gin-gonic/gin"
)

func (c *tagController) Create(ctx *gin.Context) {
	form, formError := request.NewTagRequest().ValidateJSON(ctx)
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

	tag, err := c.tagService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	infoBlocks := c.api.InfoBlock.GetAll()

	data := response.Body{
		"templates": c.templates(ctx),
		"tag":       tag,
		"collection": gin.H{
			"infoBlocks":          infoBlocks,
			"infoBlockCollection": tag.InfoBlocks,
			"relationURL":         tag.AdminURL(),
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
			c.T(ctx, "ui.success.record_created"),
		),
	)
}
