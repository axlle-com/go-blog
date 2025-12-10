package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/publisher/models"
)

type CollectionService struct {
	api *api.Api
}

func NewCollectionService(
	api *api.Api,
) *CollectionService {
	return &CollectionService{
		api: api,
	}
}

func (s *CollectionService) WithPaginate(paginator contract.Paginator, filter contract.PublisherFilter) ([]contract.Publisher, error) {
	var collection []contract.Publisher

	items, total, err := s.api.Blog.GetPublishers(paginator, filter)
	for _, item := range items {
		collection = append(collection, models.FromContract(item))
	}

	paginator.SetTotal(total)

	return collection, err
}
