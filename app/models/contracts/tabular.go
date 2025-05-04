package contracts

type Record interface {
	GetID() uint
	GetTable() string
}
