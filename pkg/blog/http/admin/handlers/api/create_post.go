package api

import (
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddPostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *controller) CreatePost(ctx *gin.Context) {
	body := AddPostRequestBody{}
	// получаем тело запроса
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post

	post.Title = body.Title
	ctx.JSON(http.StatusCreated, &post)
}
