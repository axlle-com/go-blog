package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service/mailer/queue/model"
)

func NewMailerQueueHandler(mailer contract.Mailer) contract.QueueHandler {
	return &queueHandler{
		mailer: mailer,
	}
}

type queueHandler struct {
	mailer contract.Mailer
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[mailer][queue] error: %v", err)
		return
	}

	switch action {
	case "create":
		if err := qh.create(payload); err != nil {
			logger.Errorf("[mailer][queue][create] error: %v", err)
		}
	default:
		logger.Debugf("[mailer][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) create(payload []byte) error {
	var obj model.Informer

	dec := json.NewDecoder(bytes.NewReader(payload))

	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("incorrect data format: %v", err)
	}

	return qh.mailer.SendMail(&obj)
}
