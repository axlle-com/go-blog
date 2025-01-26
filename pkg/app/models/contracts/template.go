package contracts

type Template interface {
	GetID() uint
	GetTitle() string
	GetName() string
	GetTabular() string
}
