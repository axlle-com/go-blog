package service

import (
	"github.com/axlle-com/blog/pkg/gallery/provider"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	. "github.com/axlle-com/blog/pkg/template/repository"
)

type Container interface {
	Post() PostRepository
	Category() CategoryRepository
	Template() TemplateRepository
	Request() *PostRequest
	Gallery() provider.Provider
}

type container struct {
	postRepository     PostRepository
	categoryRepository CategoryRepository
	templateRepository TemplateRepository
	gallery            provider.Provider
	request            *PostRequest
}

func NewContainer() Container {
	return &container{}
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

func (c *container) Template() TemplateRepository {
	if c.templateRepository == nil {
		c.templateRepository = NewRepo()
	}
	return c.templateRepository
}

func (c *container) Request() *PostRequest {
	if c.request == nil {
		c.request = NewPostRequest()
	}
	return c.request
}

func (c *container) Gallery() provider.Provider {
	if c.gallery == nil {
		c.gallery = provider.NewProvider()
	}
	return c.gallery
}
