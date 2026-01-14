package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/google/uuid"
)

func NewGalleryQueueHandler(
	infoBlockService *service.InfoBlockService,
	infoBlockCollectionService *service.InfoBlockCollectionService,
	infoBlockEventService *service.InfoBlockEventService,
) contract.QueueHandler {
	return &queueHandler{
		infoBlockService:           infoBlockService,
		infoBlockCollectionService: infoBlockCollectionService,
		infoBlockEventService:      infoBlockEventService,
	}
}

type queueHandler struct {
	infoBlockService           *service.InfoBlockService
	infoBlockCollectionService *service.InfoBlockCollectionService
	infoBlockEventService      *service.InfoBlockEventService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := app.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[info_block][queue] error: %v", err)
		return
	}

	switch action {
	case queue.Update:
		if err := qh.update(payload); err != nil {
			logger.Errorf("[info_block][queue][create] error: %v", err)
		}
	default:
		logger.Errorf("[info_block][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) update(payload []byte) error {
	var objects dto.Collection

	dec := json.NewDecoder(bytes.NewReader(payload))
	if err := dec.Decode(&objects); err != nil {
		return fmt.Errorf("incorrect data format: %w", err)
	}

	// Собираем уникальные resource UUID
	seen := make(map[uuid.UUID]struct{}, len(objects.ResourceBlocks))
	uuids := make([]uuid.UUID, 0, len(objects.ResourceBlocks))

	for _, rb := range objects.ResourceBlocks {
		if rb.ResourceUUID == "" {
			continue
		}

		infoBlockUUID, err := uuid.Parse(rb.ResourceUUID)
		if err != nil {
			return fmt.Errorf("invalid resource_uuid %q: %v", rb.ResourceUUID, err)
		}

		if _, ok := seen[infoBlockUUID]; ok {
			continue
		}

		seen[infoBlockUUID] = struct{}{}
		uuids = append(uuids, infoBlockUUID)
	}

	if len(uuids) == 0 {
		return nil
	}

	filter := models.NewInfoBlockFilter()
	filter.UUIDs = uuids

	return qh.infoBlockEventService.UpdatedByFilter(filter)
}
