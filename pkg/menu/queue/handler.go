package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/queue/model"
	"github.com/axlle-com/blog/pkg/menu/service"
)

func NewPublisherQueueHandler(
	menuService *service.MenuService,
	menuItemCollectionService *service.MenuItemCollectionService,
) contracts.QueueHandler {
	return &queueHandler{
		menuService:               menuService,
		menuItemCollectionService: menuItemCollectionService,
	}
}

type queueHandler struct {
	menuService               *service.MenuService
	menuItemCollectionService *service.MenuItemCollectionService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[menu][queue] error: %v", err)
		return
	}

	switch action {
	case "update":
		if err := qh.update(payload); err != nil {
			logger.Errorf("[menu][queue][update] error: %v", err)
		}
	case "delete":
		if err := qh.delete(payload); err != nil {
			logger.Errorf("[menu][queue][delete] error: %v", err)
		}
	default:
		logger.Errorf("[menu][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) update(payload []byte) error {
	obj, err := qh.parsePayload(payload)
	if err != nil {
		return err
	}

	_, err = qh.menuItemCollectionService.UpdateURLForPublisher(&obj)
	return err
}

func (qh *queueHandler) delete(payload []byte) error {
	obj, err := qh.parsePayload(payload)
	if err != nil {
		return err
	}

	_, err = qh.menuItemCollectionService.DetachPublisher(&obj)
	return err
}

func (qh *queueHandler) parsePayload(payload []byte) (obj model.Publisher, err error) {
	dec := json.NewDecoder(bytes.NewReader(payload))

	if err := dec.Decode(&obj); err != nil {
		return obj, fmt.Errorf("incorrect data format: %v", err)
	}

	return
}
