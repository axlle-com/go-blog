package service

import (
	"errors"
	app "github.com/axlle-com/blog/pkg/app/service"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/axlle-com/blog/pkg/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
)

type InfoBlockService struct {
	infoBlockRepo   repository.InfoBlockRepository
	resourceRepo    repository.InfoBlockHasResourceRepository
	galleryProvider provider.GalleryProvider
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
	galleryProvider provider.GalleryProvider,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo:   infoBlockRepo,
		resourceRepo:    resourceRepo,
		galleryProvider: galleryProvider,
	}
}

func (s *InfoBlockService) GetByID(id uint) (*InfoBlock, error) {
	return s.infoBlockRepo.GetByID(id)
}

func (s *InfoBlockService) Create(infoBlock *InfoBlock, user contracts.User) (*InfoBlock, error) {
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

func (s *InfoBlockService) Attach(resource contracts.Resource, infoBlock contracts.InfoBlock) error {
	hasRepo, err := s.resourceRepo.GetByParams(resource.GetUUID(), infoBlock.GetID())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&InfoBlockHasResource{
				ResourceUUID: resource.GetUUID(),
				InfoBlockID:  infoBlock.GetID(),
			},
		)
	}
	return nil
}

func (s *InfoBlockService) GetForResource(resource contracts.Resource) (infoBlocks []*InfoBlock, err error) {
	return s.infoBlockRepo.GetForResource(resource)
}

func (s *InfoBlockService) DeleteForResource(resource contracts.Resource) error {
	byResource, err := s.resourceRepo.GetByResource(resource)
	if err != nil {
		return err
	}
	if len(byResource) == 0 {
		return nil
	}

	// Группируем записи по InfoBlockID
	blockResources := make(map[uint][]uuid.UUID)
	for _, res := range byResource {
		blockResources[res.InfoBlockID] = append(blockResources[res.InfoBlockID], res.ResourceUUID)
	}

	var detachBlockIDs []uint
	var deleteBlockIDs []uint

	// Определяем для каждой, сколько ресурсов ей принадлежит
	for infoBlockID, resources := range blockResources {
		if len(resources) > 1 {
			detachBlockIDs = append(detachBlockIDs, infoBlockID)
		} else {
			deleteBlockIDs = append(deleteBlockIDs, infoBlockID)
		}
	}

	if len(detachBlockIDs) > 0 {
		for _, blockID := range detachBlockIDs {
			err = s.resourceRepo.DeleteByParams(resource.GetUUID(), blockID)
			if err != nil {
				return err
			}
		}
	}

	if len(deleteBlockIDs) > 0 {
		blocks, err := s.infoBlockRepo.GetByIDs(deleteBlockIDs)
		if err != nil {
			return err
		}
		err = s.DeleteInfoBlocks(blocks)
		if err != nil {
			return err
		}
	}
	return nil
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
	return s.infoBlockRepo.Delete(infoBlocks)
}

func (s *InfoBlockService) SaveFromRequest(form *BlockRequest, found *InfoBlock, user contracts.User) (infoBlock *InfoBlock, err error) {
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
		slice := make([]contracts.Gallery, 0)
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
