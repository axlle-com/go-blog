package web

import (
	"github.com/axlle-com/blog/pkg/post"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddPostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h handler) AddPost(c *gin.Context) {
	body := AddPostRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post post.Post

	post.Title = body.Title

	if result := h.DB.Create(&post); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &post)
}
