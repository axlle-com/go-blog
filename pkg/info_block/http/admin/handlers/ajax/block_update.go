package ajax

import (
	"fmt"
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
)

func (c *blockController) UpdateInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	found, err := c.blockService.FindByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	form, formError := models.NewBlockRequest().ValidateJSON(ctx)
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

	var infoBlocks []*models.InfoBlock
	if block.ID != 0 {
		var err2 error
		infoBlocks, err2 = c.blockCollectionService.GetAllForParent(block)
		if err2 != nil {
			// Если ошибка, получаем все инфоблоки
			infoBlocks, _ = c.blockCollectionService.GetAll()
		}
	} else {
		infoBlocks, _ = c.blockCollectionService.GetAll()
	}

	data := gin.H{
		"templates":  c.templates(ctx),
		"infoBlocks": infoBlocks,
		"infoBlock":  block,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":      c.RenderView("admin.info_block_inner", data, ctx),
			"url":       fmt.Sprintf("/admin/info-blocks/%d", block.ID),
			"infoBlock": block,
		},
	})
}
