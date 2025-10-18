package queue

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/file/service"
)

func NewFileQueueHandler(
	collectionService *service.CollectionService,
) contracts.QueueHandler {
	return &queueHandler{
		collectionService: collectionService,
	}
}

type queueHandler struct {
	collectionService *service.CollectionService
}

func (qh *queueHandler) Run(data []byte) {
	action, _, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[user][queue] error: %v", err)
		return
	}

	switch action {
	case "revision_received":
		if err := qh.revisionReceived(); err != nil {
			logger.Errorf("[user][queue][create] error: %v", err)
		}
	default:
		logger.Errorf("[user][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) revisionReceived() error {
	return qh.collectionService.RevisionReceived()
}
