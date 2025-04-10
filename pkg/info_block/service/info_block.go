package service

import (
	"errors"
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"gorm.io/gorm"
)

type InfoBlockService struct {
	infoBlockRepo       repository.InfoBlockRepository
	infoBlockCollection *InfoBlockCollectionService
	resourceRepo        repository.InfoBlockHasResourceRepository
	galleryProvider     provider.GalleryProvider
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	infoBlockCollection *InfoBlockCollectionService,
	resourceRepo repository.InfoBlockHasResourceRepository,
	galleryProvider provider.GalleryProvider,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo:       infoBlockRepo,
		infoBlockCollection: infoBlockCollection,
		resourceRepo:        resourceRepo,
		galleryProvider:     galleryProvider,
	}
}

func (s *InfoBlockService) GetByID(id uint) (*InfoBlock, error) {
	return s.infoBlockRepo.GetByID(id)
}

func (s *InfoBlockService) GetByIDs(ids []uint) ([]*InfoBlock, error) {
	return s.infoBlockRepo.GetByIDs(ids)
}

func (s *InfoBlockService) Create(infoBlock *InfoBlock, user contracts2.User) (*InfoBlock, error) {
	if user != nil {
		id := user.GetID()
		infoBlock.UserID = &id
	}
	if err := s.infoBlockRepo.Create(infoBlock); err != nil {
		return nil, err
	}
	return infoBlock, nil
}

func (s *InfoBlockService) Update(infoBlock *InfoBlock) (*InfoBlock, error) {
	if err := s.infoBlockRepo.Update(infoBlock); err != nil {
		return nil, err
	}

	return infoBlock, nil
}

func (s *InfoBlockService) Attach(resource contracts2.Resource, infoBlock contracts2.InfoBlock) error {
	hasRepo, err := s.resourceRepo.GetByID(infoBlock.GetRelationID())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&InfoBlockHasResource{
				ResourceUUID: resource.GetUUID(),
				InfoBlockID:  infoBlock.GetID(),
				Sort:         infoBlock.GetSort(),
				Position:     infoBlock.GetPosition(),
			},
		)
	} else {
		hasRepo.Position = infoBlock.GetPosition()
		hasRepo.Sort = infoBlock.GetSort()
		err = s.resourceRepo.Update(hasRepo)
	}

	return err
}

func (s *InfoBlockService) GetForResource(resource contracts2.Resource) []*InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResource(resource)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return s.infoBlockCollection.AggregatesResponses(infoBlocks)
}

func (s *InfoBlockService) DetachResource(resource contracts2.Resource) error {
	return s.resourceRepo.DetachResource(resource)
}

func (s *InfoBlockService) DeleteInfoBlocks(infoBlocks []*InfoBlock) (err error) {
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

func (s *InfoBlockService) Delete(infoBlocks *InfoBlock) (err error) {
	if err = s.resourceRepo.DetachInfoBlock(infoBlocks); err != nil {
		return err
	}
	return s.infoBlockRepo.Delete(infoBlocks)
}

func (s *InfoBlockService) SaveFromRequest(form *BlockRequest, found *InfoBlock, user contracts2.User) (infoBlock *InfoBlock, err error) {
	blockForm := app.LoadStruct(&InfoBlock{}, form).(*InfoBlock)

	if found == nil {
		infoBlock, err = s.Create(blockForm, user)
	} else {
		blockForm.ID = found.ID
		blockForm.UUID = found.UUID
		infoBlock, err = s.Update(blockForm)
	}

	if err != nil {
		return
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts2.Gallery, 0)
		for _, gRequest := range form.Galleries {
			if gRequest == nil {
				continue
			}

			g, err := s.galleryProvider.SaveFromForm(gRequest, infoBlock)
			if err != nil || g == nil {
				continue
			}
			slice = append(slice, g)
		}
		infoBlock.Galleries = slice
	}
	return
}
