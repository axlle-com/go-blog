package provider

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models/contract"
)

func NewPublisherProvider(
	api *api.Api,
) contract.PublisherProvider {
	return &provider{
		api: api,
	}
}

type provider struct {
	api *api.Api
}

func (p *provider) GetPublishers(paginator contract.Paginator) (collection []contract.Publisher, total int, err error) {
	posts, total, err := p.api.Blog.GetPublishers(paginator)
	if err != nil {
		return
	}

	for _, post := range posts {
		collection = append(collection, post)
	}

	return
}
