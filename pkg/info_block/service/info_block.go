package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"gorm.io/gorm"
)

type InfoBlockService struct {
	infoBlockRepo repository.InfoBlockRepository
	resourceRepo  repository.InfoBlockHasResourceRepository
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo: infoBlockRepo,
		resourceRepo:  resourceRepo,
	}
}

func (s *InfoBlockService) CreateInfoBlock(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if err := s.infoBlockRepo.Create(infoBlock); err != nil {
		return nil, err
	}
	return infoBlock, nil
}

func (s *InfoBlockService) UpdateInfoBlock(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if err := s.infoBlockRepo.Update(infoBlock); err != nil {
		return nil, err
	}

	return infoBlock, nil
}

func (s *InfoBlockService) Attach(resource contracts.Resource, infoBlock contracts.InfoBlock) error {
	hasRepo, err := s.resourceRepo.GetByParams(resource.GetUUID(), infoBlock.GetID())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.InfoBlockHasResource{
				ResourceUUID: resource.GetUUID(),
				InfoBlockID:  infoBlock.GetID(),
			},
		)
	}
	return nil
}

func (s *InfoBlockService) DeleteForResource(resource contracts.Resource) (err error) {
	byResource, err := s.resourceRepo.GetByResource(resource)
	if err != nil {
		return err
	}

	all := make(map[uint]*models.InfoBlockHasResource)
	only := make(map[uint]*models.InfoBlockHasResource)
	detach := make(map[uint]*models.InfoBlockHasResource)
	var infoBlockIDs []uint
	if byResource == nil {
		return nil
	}

	for _, r := range byResource {
		if r.ResourceUUID != resource.GetUUID() {
			all[r.InfoBlockID] = r
		} else {
			only[r.InfoBlockID] = r
		}
	}

	for id, _ := range only {
		if _, ok := all[id]; ok {
			detach[id] = all[id]
		} else {
			infoBlockIDs = append(infoBlockIDs, id)
		}
	}

	if len(detach) > 0 { // TODO need test
		for _, r := range detach {
			err = s.resourceRepo.DeleteByParams(r.ResourceUUID, r.InfoBlockID)
			if err != nil {
				return err
			}
		}
	}

	if len(infoBlockIDs) > 0 {
		infoBlocks, err := s.infoBlockRepo.GetByIDs(infoBlockIDs)
		if err != nil {
			return err
		}
		err = s.DeleteInfoBlocks(infoBlocks)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *InfoBlockService) DeleteInfoBlocks(infoBlocks []*models.InfoBlock) (err error) {
	var ids []uint
	for _, infoBlock := range infoBlocks {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.infoBlockRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}
