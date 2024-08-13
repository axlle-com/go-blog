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

func (c *controller) CreatePostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.createPost(ctx, container)
	}
}

func (c *controller) DeletePostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.DeletePost(ctx, container)
	}
}

func (c *controller) DeletePostImageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.DeletePostImage(ctx, container)
	}
}
