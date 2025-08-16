package service

import (
	"errors"
	"sync"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	tProvider "github.com/axlle-com/blog/pkg/template/provider"
	provider2 "github.com/axlle-com/blog/pkg/user/provider"
	"gorm.io/gorm"
)

type InfoBlockService struct {
	infoBlockRepo       repository.InfoBlockRepository
	infoBlockCollection *InfoBlockCollectionService
	resourceRepo        repository.InfoBlockHasResourceRepository
	galleryProvider     provider.GalleryProvider
	templateProvider    tProvider.TemplateProvider
	userProvider        provider2.UserProvider
	fileProvider        fileProvider.FileProvider
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	infoBlockCollection *InfoBlockCollectionService,
	resourceRepo repository.InfoBlockHasResourceRepository,
	galleryProvider provider.GalleryProvider,
	templateProvider tProvider.TemplateProvider,
	userProvider provider2.UserProvider,
	fileProvider fileProvider.FileProvider,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo:       infoBlockRepo,
		infoBlockCollection: infoBlockCollection,
		resourceRepo:        resourceRepo,
		galleryProvider:     galleryProvider,
		templateProvider:    templateProvider,
		userProvider:        userProvider,
		fileProvider:        fileProvider,
	}
}

func (s *InfoBlockService) GetByID(id uint) (*models.InfoBlock, error) {
	return s.infoBlockRepo.GetByID(id)
}

func (s *InfoBlockService) Aggregate(infoBlock *models.InfoBlock) *models.InfoBlock {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		if infoBlock.TemplateID != nil && *infoBlock.TemplateID != 0 {
			var err error
			infoBlock.Template, err = s.templateProvider.GetByID(*infoBlock.TemplateID)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if infoBlock.UserID != nil && *infoBlock.UserID != 0 {
			var err error
			infoBlock.User, err = s.userProvider.GetByID(*infoBlock.UserID)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		infoBlock.Galleries = s.galleryProvider.GetForResource(infoBlock)
		if err != nil {
			logger.Error(err)
		}
	}()

	wg.Wait()

	return infoBlock
}

func (s *InfoBlockService) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetByIDs(ids)
}

func (s *InfoBlockService) Create(infoBlock *models.InfoBlock, user contracts.User) (*models.InfoBlock, error) {
	if user != nil {
		id := user.GetID()
		infoBlock.UserID = &id
	}
	if err := s.infoBlockRepo.Create(infoBlock); err != nil {
		return nil, err
	}
	return infoBlock, nil
}

func (s *InfoBlockService) Update(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if err := s.infoBlockRepo.Update(infoBlock); err != nil {
		return nil, err
	}

	return infoBlock, nil
}

func (s *InfoBlockService) Attach(resource contracts.Resource, infoBlock contracts.InfoBlock) error {
	hasRepo, err := s.resourceRepo.GetByID(infoBlock.GetRelationID())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.InfoBlockHasResource{
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

func (s *InfoBlockService) GetForResource(resource contracts.Resource) []*models.InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResource(resource)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return s.infoBlockCollection.AggregatesResponses(infoBlocks)
}

func (s *InfoBlockService) DetachResource(resource contracts.Resource) error {
	return s.resourceRepo.DetachResource(resource)
}

func (s *InfoBlockService) DeleteResource(id uint) error {
	return s.resourceRepo.Delete(id)
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

func (s *InfoBlockService) Delete(infoBlocks *models.InfoBlock) (err error) {
	if err = s.resourceRepo.DetachInfoBlock(infoBlocks); err != nil {
		return err
	}
	return s.infoBlockRepo.Delete(infoBlocks)
}

func (s *InfoBlockService) SaveFromRequest(form *models.BlockRequest, found *models.InfoBlock, user contracts.User) (infoBlock *models.InfoBlock, err error) {
	blockForm := app.LoadStruct(&models.InfoBlock{}, form).(*models.InfoBlock)

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

			g, err := s.galleryProvider.SaveForm(gRequest, infoBlock)
			if err != nil || g == nil {
				continue
			}
			slice = append(slice, g)
		}
		infoBlock.Galleries = slice
	}
	return
}

func (s *InfoBlockService) DeleteImageFile(block *models.InfoBlock) error {
	if block.Image == nil {
		return nil
	}
	err := s.fileProvider.DeleteFile(*block.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		logger.Errorf("[DeleteImageFile] Error: %v", err)
	}
	block.Image = nil
	return nil
}
