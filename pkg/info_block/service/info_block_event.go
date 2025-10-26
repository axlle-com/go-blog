package service

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/queue/job"
)

type InfoBlockEventService struct {
	queue contracts.Queue
}

func NewInfoBlockEventService(
	queue contracts.Queue,
) *InfoBlockEventService {
	return &InfoBlockEventService{
		queue: queue,
	}
}
func (s *InfoBlockEventService) StartJob(collection []*models.InfoBlockResponse) {
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

	newJob := job.NewInfoBlockJob(
		&dto.Collection{ResourceBlocks: slice},
		"update",
	)

	s.queue.Enqueue(newJob, 0)
}
