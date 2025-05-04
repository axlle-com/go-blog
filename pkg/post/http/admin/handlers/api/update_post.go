package api

import (
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdatePostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *controller) UpdatePost(ctx *gin.Context) {
	body := UpdatePostRequestBody{}
	// получаем тело запроса
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post
	post.Title = body.Title

	ctx.JSON(http.StatusOK, &post)
}
