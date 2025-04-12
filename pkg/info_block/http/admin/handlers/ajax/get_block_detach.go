package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) DetachInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	err := c.blockService.DeleteResource(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			nil,
			"Запись удалена",
			nil,
		),
	)
}
