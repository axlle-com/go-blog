package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) GetInfoBlockCard(ctx *gin.Context) {
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
	found = c.blockService.Aggregate(found)

	data := response.Body{
		"infoBlock": found,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view":      c.RenderView("admin.info_block_card", data, ctx),
				"url":       nil,
				"infoBlock": found,
			},
			"Запись найдена",
			nil,
		),
	)
}
