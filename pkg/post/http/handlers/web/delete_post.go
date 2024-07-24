package web

import (
	"github.com/axlle-com/blog/pkg/post/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	postRepo := repository.NewRepository()
	num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Ошибка преобразования:", err)
		return
	}

	if err := postRepo.DeletePost(uint(num)); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}
