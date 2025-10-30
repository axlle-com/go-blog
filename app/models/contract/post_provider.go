package contract

type PostProvider interface {
	GetPosts() []Post
	GetPublishers() ([]Publisher, error)
}
