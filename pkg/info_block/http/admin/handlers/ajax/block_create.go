package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/app/http/response"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) CreateInfoBlock(ctx *gin.Context) {
	form, formError := NewBlockRequest().ValidateJSON(ctx)
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

	block, err := c.blockService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	templates := c.templateProvider.GetAll()

	data := response.Body{
		"templates": templates,
		"infoBlock": block,
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view":      c.RenderView("admin.info_block_inner", data, ctx),
				"url":       fmt.Sprintf("/admin/info-blocks/%d", block.ID),
				"infoBlock": block,
			},
			"Запись создана",
		),
	)
}
