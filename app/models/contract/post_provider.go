package contract

type BlogProvider interface {
	GetPosts() []Post
	GetPublishers(Paginator, PublisherFilter) ([]Publisher, int, error)
}
