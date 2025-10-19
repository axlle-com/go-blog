package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/queue/model"
	"github.com/axlle-com/blog/pkg/message/service"
)

func NewMessageQueueHandler(
	messageService *service.MessageService,
	messageCollectionService *service.MessageCollectionService,
) contracts.QueueHandler {
	return &queueHandler{
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
	}
}

type queueHandler struct {
	messageService           *service.MessageService
	messageCollectionService *service.MessageCollectionService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[user][queue] error: %v", err)
		return
	}

	switch action {
	case "create":
		if err := qh.create(payload); err != nil {
			logger.Errorf("[message][queue][create] error: %v", err)
		}
	default:
		logger.Errorf("[message][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) create(payload []byte) error {
	var obj model.Message

	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("incorrect data format: %v", err)
	}

	_, err := qh.messageService.Create(obj.Model(), obj.UserUUID)
	return err
}
