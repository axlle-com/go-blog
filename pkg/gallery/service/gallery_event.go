package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/queue/job"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

type GalleryEvent struct {
	queue        contract.Queue
	imageService *ImageService
	galleryRepo  repository.GalleryRepository
	resourceRepo repository.GalleryResourceRepository
}

func NewGalleryEvent(
	queue contract.Queue,
	imageService *ImageService,
	galleryRepo repository.GalleryRepository,
	resource repository.GalleryResourceRepository,
) *GalleryEvent {
	return &GalleryEvent{
		queue:        queue,
		imageService: imageService,
		galleryRepo:  galleryRepo,
		resourceRepo: resource,
	}
}

func (e *GalleryEvent) DeletingGallery(g *models.Gallery) (err error) {
	has, _ := e.resourceRepo.GetByGalleryID(g.ID) // @todo сразу delete
	if has != nil {
		if err = e.resourceRepo.Delete(g.ID); err != nil {
			return err
		}
	}

	err = e.imageService.DeleteImages(g.Images)
	if err != nil {
		return err
	}

	return nil
}

func (e *GalleryEvent) DeletedGallery(g *models.Gallery) (err error) {
	return err
}

func (e *GalleryEvent) UpdateTrigger(ids []uint) {
	filter := models.NewGalleryFilter()
	filter.IDs = ids
	collection, err := e.galleryRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[gallery][GalleryEvent][UpdateTrigger] error: %v", err)
		return
	}

	if len(collection) == 0 {
		return
	}

	// Use map to deduplicate by ResourceUUID + BlockUUID combination
	uniqueBlocks := make(map[string]*dto.ResourceBlock)

	for _, response := range collection {
		resourceUUID := response.ResourceUUID.String()
		blockUUID := response.UUID.String()

		// Create unique key for deduplication
		key := resourceUUID + ":" + blockUUID

		// Only add if not already present
		if _, exists := uniqueBlocks[key]; !exists {
			uniqueBlocks[key] = &dto.ResourceBlock{
				ResourceUUID: resourceUUID,
				BlockUUID:    blockUUID,
			}
		}
	}

	// Convert map values to slice
	var slice []*dto.ResourceBlock
	for _, block := range uniqueBlocks {
		slice = append(slice, block)
	}

	newJob := job.NewGalleryJob(
		&dto.Collection{ResourceBlocks: slice},
		queue.Update,
	)

	e.queue.Enqueue(newJob, 0)

	return
}
