package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/blog/service"
)

func NewPostProvider(
	postRepo repository.PostRepository,
	postService *service.PostService,
) contracts.PostProvider {
	return &provider{
		postRepo:    postRepo,
		postService: postService,
	}
}

type provider struct {
	postRepo    repository.PostRepository
	postService *service.PostService
}

func (p *provider) GetAll() []contracts.Post {
	var collection []contracts.Post
	galleries, err := p.postRepo.GetAll()
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}
