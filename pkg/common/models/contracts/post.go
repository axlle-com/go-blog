package contracts

type Resource interface {
	GetID() uint
	GetResource() string
}
