package web

import (
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdatePostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h handler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	body := UpdatePostRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post Post

	if result := h.DB.First(&post, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	post.Title = body.Title

	h.DB.Save(&post)

	c.JSON(http.StatusOK, &post)
}
