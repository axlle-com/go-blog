package ajax

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (c *controller) AddPostInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.postService.GetByID(id)
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

	postInfoBlocks, err := c.infoBlock.Attach(uint(infoBlockId), found)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	infoBlocks := c.infoBlock.GetAll()
	data := gin.H{
		"infoBlocks":     infoBlocks,
		"postInfoBlocks": postInfoBlocks,
		"relationID":     found.ID,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":           c.RenderView("admin.info_block_in_item_inner", data, ctx),
			"infoBlocks":     infoBlocks,
			"postInfoBlocks": postInfoBlocks,
		},
	})
}
