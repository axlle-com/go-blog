package web

import (
	"github.com/axlle-com/blog/pkg/post/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *controller) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	postRepo := repository.NewPostRepository()
	num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Ошибка преобразования:", err)
		return
	}

	if err := postRepo.Delete(uint(num)); err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.Status(http.StatusOK)
}
