package api

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddPostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func AddPost(c *gin.Context) {
	body := AddPostRequestBody{}
	h := db.GetDB()
	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post

	post.Title = body.Title

	if result := h.Create(&post); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &post)
}
