package api

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdatePostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *controller) UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	body := UpdatePostRequestBody{}
	h := db.GetDB()
	// получаем тело запроса
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post models.Post

	if result := h.First(&post, id); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	post.Title = body.Title

	h.Save(&post)

	ctx.JSON(http.StatusOK, &post)
}
