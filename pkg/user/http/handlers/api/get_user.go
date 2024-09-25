package api

import (
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetUser(c *gin.Context) {
	id := c.Param("id")
	value, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Println("Ошибка преобразования:", err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	result, err := repository.NewRepo().GetByID(uint(value))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
