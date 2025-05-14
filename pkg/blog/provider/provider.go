package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/blog/service"
)

type PostProvider interface {
	GetAll() []contracts.Post
}

func NewPostProvider(
	postRepo repository.PostRepository,
	service *service.PostService,
) PostProvider {
	return &provider{
		postRepo: postRepo,
		service:  service,
	}
}

type provider struct {
	postRepo repository.PostRepository
	service  *service.PostService
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
