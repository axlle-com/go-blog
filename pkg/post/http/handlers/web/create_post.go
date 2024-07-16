package web

import (
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h handler) CreatePost(c *gin.Context) {
	body := models.CreatePostRequest{}
	postRepo := repository.NewPostRepository(h.DB)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post Post

	post.Title = body.Title

	if err := postRepo.CreatePost(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post})
}
