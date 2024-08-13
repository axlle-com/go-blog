package ajax

import (
	"bytes"
	common "github.com/axlle-com/blog/pkg/common/models"
	handlers "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePost(ctx *gin.Context, ctr Container) {
	id := c.getID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}
	postRepo := NewPostRepo()
	if err := postRepo.Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	paginator := common.NewPaginator(ctx)
	body := handlers.NewResponse(paginator).GetForWeb()

	var buf bytes.Buffer
	originalWriter := ctx.Writer

	wrappedWriter := &ResponseWriterWrapper{
		ResponseWriter: ctx.Writer,
		Buffer:         &buf,
	}
	ctx.Writer = wrappedWriter
	ctx.HTML(http.StatusOK, "admin.posts_inner", body)
	ctx.Writer = originalWriter
	body["view"] = buf.String()
	ctx.JSON(http.StatusOK, gin.H{"data": body})
}
