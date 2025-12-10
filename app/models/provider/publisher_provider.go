package provider

import "github.com/axlle-com/blog/app/models/contract"

type PublisherProvider interface {
	GetPublishers(contract.Paginator, contract.PublisherFilter) ([]contract.Publisher, int, error)
}
