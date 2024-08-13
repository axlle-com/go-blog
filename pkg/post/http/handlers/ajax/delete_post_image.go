package ajax

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	. "github.com/axlle-com/blog/pkg/post/models"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePostImage(ctx *gin.Context, ctr Container) {
	id := c.getID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}
	postRepo := NewPostRepo()
	post, err := postRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	post.DeleteImageFile()
	err = postRepo.Update(post)
	if err != nil {
		logger.New().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
	return
}
