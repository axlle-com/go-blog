package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/service"
)

func NewPostProvider(
	postService *service.PostService,
	postCollectionService *service.PostCollectionService,
	categoriesService *service.CategoriesService,
	tagCollectionService *service.TagCollectionService,
) contracts.PostProvider {
	return &provider{
		postService:           postService,
		postCollectionService: postCollectionService,
		categoriesService:     categoriesService,
		tagCollectionService:  tagCollectionService,
	}
}

type provider struct {
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoriesService     *service.CategoriesService
	tagCollectionService  *service.TagCollectionService
}

func (p *provider) GetPosts() []contracts.Post {
	var collection []contracts.Post
	posts, err := p.postCollectionService.GetAll()
	if err == nil {
		for _, post := range posts {
			collection = append(collection, post)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetPublishers() ([]contracts.Publisher, error) { // @todo paginate!!!
	var collection []contracts.Publisher

	posts, err := p.postCollectionService.GetAll()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		collection = append(collection, post)
	}

	categories, err := p.categoriesService.GetAll()
	if err != nil {
		return collection, err
	}

	for _, category := range categories {
		collection = append(collection, category)
	}

	tags, err := p.tagCollectionService.GetAll()
	if err != nil {
		return collection, err
	}

	for _, tag := range tags {
		collection = append(collection, tag)
	}

	return collection, err
}
