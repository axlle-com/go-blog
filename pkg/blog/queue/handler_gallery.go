package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func NewGalleryQueueHandler(
	categoriesService *service.CategoriesService,
	postCollectionService *service.PostCollectionService,
	tagCollectionService *service.TagCollectionService,
	api *api.Api,
) contract.QueueHandler {
	return &queueGalleryHandler{
		categoriesService:     categoriesService,
		postCollectionService: postCollectionService,
		tagCollectionService:  tagCollectionService,
		api:                   api,
	}
}

type queueGalleryHandler struct {
	categoriesService     *service.CategoriesService
	postCollectionService *service.PostCollectionService
	tagCollectionService  *service.TagCollectionService
	api                   *api.Api
}

func (qh *queueGalleryHandler) Run(data []byte) {
	action, payload, err := app.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[blog][queue] error: %v", err)
		return
	}

	switch action {
	case "update":
		if err := qh.update(payload); err != nil {
			logger.Errorf("[blog][queue][create] error: %v", err)
		}
	default:
		logger.Errorf("[blog][queue] unknown action: %s", action)
	}
}

func (qh *queueGalleryHandler) update(payload []byte) error {
	var objects dto.Collection

	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&objects); err != nil {
		return fmt.Errorf("incorrect data format: %w", err)
	}

	type bundle struct {
		uuids    []uuid.UUID
		snapshot []dto.Gallery
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
			logger.Errorf("[blog][update] invalid resource_uuid %q: %v", resUUID, err)
			continue
		}

		blocks := qh.api.Gallery.GetForResourceUUID(resUUID)

		byRes[resUUID] = &bundle{
			uuids:    []uuid.UUID{newUUID},
			snapshot: dto.MapGalleries(blocks),
		}
	}

	agg := errutil.New()

	for _, b := range byRes {
		raw, err := json.Marshal(b.snapshot)
		if err != nil {
			agg.Add(fmt.Errorf("marshal galleries_snapshot: %w", err))
			continue
		}

		patch := map[string]any{
			"galleries_snapshot": datatypes.JSON(raw),
		}

		if _, err := qh.categoriesService.UpdateFieldsByUUIDs(b.uuids, patch); err != nil {
			agg.Add(err)
		}
		if _, err := qh.postCollectionService.UpdateFieldsByUUIDs(b.uuids, patch); err != nil {
			agg.Add(err)
		}
		if _, err := qh.tagCollectionService.UpdateFieldsByUUIDs(b.uuids, patch); err != nil {
			agg.Add(err)
		}
	}

	return agg.Error()
}
