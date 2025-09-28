package service

import (
	"encoding/json"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
)

type informer struct {
	from    string
	to      string
	subject string
	body    string
}

func NewInformer(body, subject string) contracts.MailRequest {
	return &informer{
		body:    body,
		subject: subject,
		from:    config.Config().SMTPUsername(),
		to:      "axlle@mail.ru",
	}
}

func (i *informer) From() string {
	return i.from
}

func (i *informer) To() string {
	return i.to
}

func (i *informer) Subject() string {
	return i.subject
}

func (i *informer) Body() string {
	return i.body
}

func (i *informer) ToString() string {
	raw, err := json.Marshal(map[string]string{
		"from":    i.from,
		"to":      i.to,
		"subject": i.subject,
		"body":    i.body,
	})
	if err != nil {
		logger.Errorf("[informer][ToString]Error: %v", err)
	}
	return string(raw)
}
