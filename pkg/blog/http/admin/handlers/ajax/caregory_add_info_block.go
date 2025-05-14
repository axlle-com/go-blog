package ajax

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (c *categoryController) AddPostInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	idParam := ctx.Param("info_block_id")
	infoBlockId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	ifoBlockCollection, err := c.infoBlockProvider.Attach(uint(infoBlockId), found)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	infoBlocks := c.infoBlockProvider.GetAll()
	data := gin.H{
		"infoBlocks":         infoBlocks,
		"ifoBlockCollection": ifoBlockCollection,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":               c.RenderView("admin.info_block_in_item_inner", data, ctx),
			"infoBlocks":         infoBlocks,
			"ifoBlockCollection": ifoBlockCollection,
		},
	})
}
