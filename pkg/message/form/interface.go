package form

import "github.com/axlle-com/blog/pkg/message/models"

type Form interface {
	Data() string
	Name() string
	Title() string
	GetUserUUID() string
	GetFrom() string
	Model() *models.Message
}
