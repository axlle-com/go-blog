package contract

type Template interface {
	GetID() uint
	GetTitle() string
	GetName() string
	GetFullName(resourceName string) string
	GetResourceName() string
}
