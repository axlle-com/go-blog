package ajax

import (
	"net/http"
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func (c *categoryController) AddPostInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	found, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	idParam := ctx.Param("info_block_id")
	infoBlockId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	infoBlockCollection, err := c.infoBlockProvider.Attach(uint(infoBlockId), found.UUID.String())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	infoBlocks := c.infoBlockProvider.GetAll()
	data := gin.H{
		"infoBlocks":          infoBlocks,
		"infoBlockCollection": infoBlockCollection,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":                c.RenderView("admin.info_block_in_item_inner", data, ctx),
			"infoBlocks":          infoBlocks,
			"infoBlockCollection": infoBlockCollection,
		},
	})
}
