package contracts

type PostProvider interface {
	GetPosts() []Post
	GetPublishers() ([]Publisher, error)
}
