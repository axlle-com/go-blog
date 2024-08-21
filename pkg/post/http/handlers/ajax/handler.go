package ajax

import (
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
)

func (c *controller) UpdatePostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.updatePost(ctx, container)
	}
}

func (c *controller) DeletePostImageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.deletePostImage(ctx, container)
	}
}

func (c *controller) FilterPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.filterPosts(ctx, container)
	}
}
