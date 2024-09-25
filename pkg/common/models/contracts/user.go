package contracts

type User interface {
	GetID() uint
	GetFirstName() string
	GetLastName() string
	GetPatronymic() string
	GetPhone() string
	GetEmail() string
	GetStatus() int8
	GetRoles() int8
	GetPermissions() int8
}
