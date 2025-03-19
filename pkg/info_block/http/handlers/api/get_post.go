package api

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	h := db.GetDB()
	var post models.Post

	if result := h.First(&post, id); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, &post)
}
