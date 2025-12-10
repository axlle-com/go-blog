package provider

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models/contract"
	appProvider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/publisher/models"
)

func NewPublisherProvider(api *api.Api) appProvider.PublisherProvider {
	return &provider{
		api: api,
	}
}

type provider struct {
	api *api.Api
}

func (p *provider) GetPublishers(paginator contract.Paginator, filter contract.PublisherFilter) (collection []contract.Publisher, total int, err error) {
	posts, total, err := p.api.Blog.GetPublishers(paginator, filter)
	if err != nil {
		return
	}

	for _, item := range posts {
		collection = append(collection, models.FromContract(item))
	}

	return
}
