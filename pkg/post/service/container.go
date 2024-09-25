package service

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	. "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type Container interface {
	User() user.User
	Post() PostRepository
	Category() CategoryRepository
	Template() Template
	Request() *PostRequest
	Gallery() provider.Gallery
	Paginator(ctx *gin.Context) contracts.Paginator
	Filter() *PostFilter
}

type container struct {
	user               user.User
	postRepository     PostRepository
	categoryRepository CategoryRepository
	template           Template
	gallery            provider.Gallery
	request            *PostRequest
	paginator          contracts.Paginator
	filter             *PostFilter
}

func NewContainer() Container {
	return &container{}
}

func (c *container) User() user.User {
	if c.user == nil {
		c.user = user.Provider()
	}
	return c.user
}

func (c *container) Post() PostRepository {
	if c.postRepository == nil {
		c.postRepository = NewPostRepo()
	}
	return c.postRepository
}

func (c *container) Category() CategoryRepository {
	if c.categoryRepository == nil {
		c.categoryRepository = NewCategoryRepo()
	}
	return c.categoryRepository
}

func (c *container) Template() Template {
	if c.template == nil {
		c.template = Provider()
	}
	return c.template
}

func (c *container) Request() *PostRequest {
	if c.request == nil {
		c.request = NewPostRequest()
	}
	return c.request
}

func (c *container) Gallery() provider.Gallery {
	if c.gallery == nil {
		c.gallery = provider.Provider()
	}
	return c.gallery
}

func (c *container) Paginator(ctx *gin.Context) contracts.Paginator {
	if c.paginator == nil {
		c.paginator = models.Paginator(ctx.Request.URL.Query())
	}
	return c.paginator
}

func (c *container) Filter() *PostFilter {
	if c.filter == nil {
		c.filter = NewPostFilter()
	}
	return c.filter
}
