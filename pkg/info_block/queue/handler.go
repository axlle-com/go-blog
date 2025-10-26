package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/google/uuid"
)

func NewGalleryQueueHandler(
	infoBlockService *service.InfoBlockService,
	infoBlockEventService *service.InfoBlockEventService,
) contracts.QueueHandler {
	return &queueHandler{
		infoBlockService:      infoBlockService,
		infoBlockEventService: infoBlockEventService,
	}
}

type queueHandler struct {
	infoBlockService      *service.InfoBlockService
	infoBlockEventService *service.InfoBlockEventService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := app.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[info_block][queue] error: %v", err)
		return
	}

	switch action {
	case "update":
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
	dec.DisallowUnknownFields()
	if err := dec.Decode(&objects); err != nil {
		return fmt.Errorf("incorrect data format: %w", err)
	}

	type bundle struct {
		uuids    []uuid.UUID
		snapshot []dto.InfoBlock
	}
	byRes := make(map[string]*bundle)

	for _, ib := range objects.ResourceBlocks {
		resUUID := ib.ResourceUUID
		if resUUID == "" {
			continue
		}

		if _, seen := byRes[resUUID]; seen {
			continue
		}

		newUUID, err := uuid.Parse(resUUID)
		if err != nil {
			return fmt.Errorf("invalid resource_uuid %q: %v", resUUID, err)
		}

		filter := models.NewInfoBlockFilter()
		filter.UUIDs = []uuid.UUID{newUUID}

		qh.infoBlockEventService.StartJob(qh.infoBlockService.GetForResourceByFilter(filter))
	}

	return nil
}
