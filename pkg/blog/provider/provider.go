package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/service"
)

func NewBlogProvider(
	postService *service.PostService,
	postCollectionService *service.PostCollectionService,
	categoriesService *service.CategoriesService,
	tagCollectionService *service.TagCollectionService,
) contract.BlogProvider {
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

func (p *provider) GetPosts() []contract.Post {
	var collection []contract.Post
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

func (p *provider) GetPublishers(paginator contract.Paginator) (collection []contract.Publisher, total int, err error) {
	newPaginator := paginator.Clone()
	total = newPaginator.GetTotal()

	posts, err := p.postCollectionService.WithPaginate(paginator, nil)
	if err != nil {
		return
	}

	if total <= paginator.GetTotal() {
		total = paginator.GetTotal()
	}

	for _, post := range posts {
		collection = append(collection, post)
	}

	categories, err := p.categoriesService.WithPaginate(paginator, nil)
	if err != nil {
		return
	}

	if total <= paginator.GetTotal() {
		total = paginator.GetTotal()
	}

	for _, category := range categories {
		collection = append(collection, category)
	}

	tags, err := p.tagCollectionService.WithPaginate(paginator, nil)
	if err != nil {
		return
	}

	if total <= paginator.GetTotal() {
		total = paginator.GetTotal()
	}

	for _, tag := range tags {
		collection = append(collection, tag)
	}

	return
}
