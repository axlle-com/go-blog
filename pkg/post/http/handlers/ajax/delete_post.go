package ajax

import (
	common "github.com/axlle-com/blog/pkg/common/models"
	handlers "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) deletePost(ctx *gin.Context, ctr Container) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}
	if err := ctr.Post().Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	paginator := common.NewPaginator(ctx)
	data := handlers.NewResponse(paginator).GetForAjax()

	data["view"] = c.RenderView("admin.posts_inner", data, ctx)
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
