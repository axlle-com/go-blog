package api

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddPostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *controller) CreatePost(ctx *gin.Context) {
	body := AddPostRequestBody{}
	h := db.GetDB()
	// получаем тело запроса
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post

	post.Title = body.Title

	if result := h.Create(&post); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, &post)
}
