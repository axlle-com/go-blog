package provider

import "github.com/axlle-com/blog/app/models/contract"

type BlogProvider interface {
	GetPosts() []contract.Post
	GetPublishers(contract.Paginator, contract.PublisherFilter) ([]contract.Publisher, int, error)
}
