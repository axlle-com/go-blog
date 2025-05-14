package ajax

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) GetInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	found, err := c.blockService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	ctx.JSON(
		http.StatusCreated,
		response.OK(
			response.Body{
				"view":      nil,
				"url":       nil,
				"infoBlock": found,
			},
			"Запись создана",
			nil,
		),
	)
}
