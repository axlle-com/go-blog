package contract

type PublisherProvider interface {
	GetPublishers(Paginator) ([]Publisher, int, error)
}
