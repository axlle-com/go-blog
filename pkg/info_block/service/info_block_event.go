package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/queue/job"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
)

type InfoBlockEventService struct {
	queue         contract.Queue
	infoBlockRepo repository.InfoBlockRepository
}

func NewInfoBlockEventService(
	queue contract.Queue,
	infoBlockRepo repository.InfoBlockRepository,
) *InfoBlockEventService {
	return &InfoBlockEventService{
		queue:         queue,
		infoBlockRepo: infoBlockRepo,
	}
}

func (s *InfoBlockEventService) Created(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	var ids []uint
	var err error

	if infoBlock.InfoBlockID != nil {
		ids, err = PathToAncestorIDs(infoBlock.PathLtree)
		if err != nil {
			return infoBlock, err
		}
	}

	filter := models.NewInfoBlockFilter()
	filter.IDs = ids

	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return infoBlock, err
	}

	s.startJob(infoBlocks, queue.Create)

	return infoBlock, err
}

func (s *InfoBlockEventService) Updated(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	var ids []uint
	var err error

	if infoBlock.InfoBlockID != nil {
		ids, err = PathToAncestorIDs(infoBlock.PathLtree)
		if err != nil {
			return infoBlock, err
		}
	}

	filter := models.NewInfoBlockFilter()
	filter.IDs = append(ids, infoBlock.ID)

	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return infoBlock, err
	}

	s.startJob(infoBlocks, queue.Update)

	return infoBlock, err
}

func (s *InfoBlockEventService) UpdatedByFilter(filter *models.InfoBlockFilter) error {
	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return err
	}

	s.startJob(infoBlocks, queue.Update)

	return err
}

func (s *InfoBlockEventService) Deleted(infoBlock *models.InfoBlock) error {
	var ids []uint
	var err error

	if infoBlock.InfoBlockID != nil {
		ids, err = PathToAncestorIDs(infoBlock.PathLtree)
		if err != nil {
			return err
		}
	}

	filter := models.NewInfoBlockFilter()
	filter.IDs = append(ids, infoBlock.ID)

	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return err
	}

	s.startJob(infoBlocks, queue.Delete)

	return err
}

func (s *InfoBlockEventService) DeletedByFilter(filter *models.InfoBlockFilter) error {
	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return err
	}

	s.startJob(collection, queue.Delete)

	return err
}

func (s *InfoBlockEventService) Attached(resourceUUID uuid.UUID) error {
	filter := models.NewInfoBlockFilter()
	filter.ResourceUUID = &resourceUUID

	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		return err
	}

	s.startJob(infoBlocks, queue.Attach)

	return err
}

func (s *InfoBlockEventService) startJob(collection []*models.InfoBlockResponse, action string) {
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
		action,
	)

	s.queue.Enqueue(newJob, 0)
}

func PathToAncestorIDs(path string) ([]uint, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return []uint{}, nil
	}

	parts := strings.Split(path, ".")
	ids := make([]uint, 0, len(parts))

	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad path segment %q in %q: %w", s, path, err)
		}
		ids = append(ids, uint(v))
	}

	// убираем последний сегмент (это ID самого узла)
	if len(ids) > 0 {
		ids = ids[:len(ids)-1]
	}

	return ids, nil
}
