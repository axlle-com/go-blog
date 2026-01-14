package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/pkg/analytic/queue/model"
	"github.com/axlle-com/blog/pkg/analytic/service"
)

func NewAnalyticQueueHandler(
	analyticService *service.AnalyticService,
	analyticCollectionService *service.AnalyticCollectionService,
) contract.QueueHandler {
	return &queueHandler{
		analyticService:           analyticService,
		analyticCollectionService: analyticCollectionService,
	}
}

type queueHandler struct {
	analyticService           *service.AnalyticService
	analyticCollectionService *service.AnalyticCollectionService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[analytic][queue] error: %v", err)
		return
	}

	switch action {
	case queue.Create:
		if err := qh.create(payload); err != nil {
			logger.Errorf("[analytic][queue][create] error: %v", err)
		}
	default:
		logger.Debugf("[analytic][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) create(payload []byte) error {
	var obj model.Analytic

	dec := json.NewDecoder(bytes.NewReader(payload))

	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("incorrect data format: %v", err)
	}

	_, err := qh.analyticService.Create(obj.Model())
	return err
}
