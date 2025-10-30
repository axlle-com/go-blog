package contract

type Template interface {
	GetID() uint
	GetTitle() string
	GetName() string
	GetResourceName() string
	GetThemeName() string
}
