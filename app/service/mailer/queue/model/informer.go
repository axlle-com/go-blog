package model

import (
	"encoding/json"

	"github.com/axlle-com/blog/app/logger"
)

type Informer struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (i *Informer) GetFrom() string {
	return i.From
}

func (i *Informer) GetTo() string {
	return i.To
}

func (i *Informer) GetSubject() string {
	return i.Subject
}

func (i *Informer) GetBody() string {
	return i.Body
}

func (i *Informer) ToString() string {
	raw, err := json.Marshal(map[string]string{
		"from":    i.From,
		"to":      i.To,
		"subject": i.Subject,
		"body":    i.Body,
	})
	if err != nil {
		logger.Errorf("[Informer][ToString]Error: %v", err)
	}
	return string(raw)
}
