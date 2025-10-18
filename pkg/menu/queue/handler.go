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
) contracts.QueueHandler {
	return &queueHandler{
		menuService: menuService,
	}
}

type queueHandler struct {
	menuService *service.MenuService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[menu][queue] error: %v", err)
		return
	}

	switch action {
	case "update":
		if err := qh.create(payload); err != nil {
			logger.Errorf("[menu][queue][update] error: %v", err)
		}
	default:
		logger.Errorf("[menu][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) create(payload []byte) error {
	var obj model.Publisher

	dec := json.NewDecoder(bytes.NewReader(payload))

	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("incorrect data format: %v", err)
	}

	logger.Dump(obj)
	//_, err := qh.menuService.Create(obj.Model(), "")
	return nil
}
