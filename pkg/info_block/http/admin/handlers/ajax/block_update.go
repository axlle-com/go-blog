package ajax

import (
	"fmt"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) UpdateInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.blockService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form, formError := NewBlockRequest().ValidateJSON(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors":  formError.Errors,
				"message": formError.Message,
			})
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	block, err := c.blockService.SaveFromRequest(form, found, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	templates := c.templateProvider.GetAll()

	data := gin.H{
		"templates": templates,
		"infoBlock": block,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":      c.RenderView("admin.block_inner", data, ctx),
			"url":       fmt.Sprintf("/admin/info-blocks/%d", block.ID),
			"infoBlock": block,
		},
	})
}
