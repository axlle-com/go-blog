package contracts

type PostProvider interface {
	GetAll() []Post
}
