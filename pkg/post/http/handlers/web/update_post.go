package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	body := UpdatePostRequestBody{}
	h := db.GetDB()
	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var post Post

	if result := h.First(&post, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	post.Title = body.Title

	h.Save(&post)

	c.JSON(http.StatusOK, &post)
}
