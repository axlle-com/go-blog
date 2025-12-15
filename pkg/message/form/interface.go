package form

import "github.com/axlle-com/blog/pkg/message/models"

type Form interface {
	Data() string
	Name() string
	Title() string
	GetUserUUID() string
	GetUserName() *string
	GetFrom() string
	GetTo() string
	Model() *models.Message
}
