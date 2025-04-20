package contracts

import (
	"github.com/axlle-com/blog/pkg/message/form"
)

type MailService interface {
	SendContact(form form.Form)
}
