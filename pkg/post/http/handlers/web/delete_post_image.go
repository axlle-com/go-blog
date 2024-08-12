package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePostImage(ctx *gin.Context) {
	id := c.getID(ctx)
	if id == 0 {
		return
	}
	postRepo := repository.NewPostRepository()
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
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
	return
}
