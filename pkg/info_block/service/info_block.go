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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InfoBlockService struct {
	infoBlockRepo         repository.InfoBlockRepository
	infoBlockCollection   *InfoBlockCollectionService
	resourceRepo          repository.InfoBlockHasResourceRepository
	infoBlockEventService *InfoBlockEventService
	galleryProvider       provider.GalleryProvider
	templateProvider      tProvider.TemplateProvider
	userProvider          provider2.UserProvider
	fileProvider          fileProvider.FileProvider
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	infoBlockCollection *InfoBlockCollectionService,
	resourceRepo repository.InfoBlockHasResourceRepository,
	infoBlockEventService *InfoBlockEventService,
	galleryProvider provider.GalleryProvider,
	templateProvider tProvider.TemplateProvider,
	userProvider provider2.UserProvider,
	fileProvider fileProvider.FileProvider,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo:         infoBlockRepo,
		infoBlockCollection:   infoBlockCollection,
		resourceRepo:          resourceRepo,
		infoBlockEventService: infoBlockEventService,
		galleryProvider:       galleryProvider,
		templateProvider:      templateProvider,
		userProvider:          userProvider,
		fileProvider:          fileProvider,
	}
}

func (s *InfoBlockService) FindByID(id uint) (*models.InfoBlock, error) {
	return s.infoBlockRepo.FindByID(id)
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

func (s *InfoBlockService) GetForResourceByFilter(filter *models.InfoBlockFilter) []*models.InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return s.infoBlockCollection.AggregatesResponses(infoBlocks)
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

	filter := models.NewInfoBlockFilter()
	filter.ID = &infoBlock.ID

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][Update] error: %v", err)
	}

	s.infoBlockEventService.StartJob(collection)

	return infoBlock, nil
}

func (s *InfoBlockService) Attach(resourceUUID uuid.UUID, infoBlock contracts.InfoBlock) error {
	hasRepo, err := s.resourceRepo.FindByID(infoBlock.GetRelationID())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.InfoBlockHasResource{
				ResourceUUID: resourceUUID,
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

	if err != nil {
		return err
	}

	filter := models.NewInfoBlockFilter()
	filter.ResourceUUID = &resourceUUID

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][Attach] error: %v", err)
	}

	s.infoBlockEventService.StartJob(collection)

	return nil
}

func (s *InfoBlockService) DeleteByResourceUUID(resourceUUID uuid.UUID) error {
	filter := models.NewInfoBlockFilter()
	filter.ResourceUUID = &resourceUUID

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][Attach] error: %v", err)
	}

	err = s.resourceRepo.DeleteByResourceUUID(resourceUUID)
	if err != nil {
		return err
	}

	s.infoBlockEventService.StartJob(collection)

	return nil
}

func (s *InfoBlockService) DeleteHasResourceByID(id uint) error {
	filter := models.NewInfoBlockFilter()
	filter.RelationID = &id

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][DeleteHasResourceByID] error: %v", err)
	}

	if err := s.resourceRepo.DeleteByID(id); err == nil {
		s.infoBlockEventService.StartJob(collection)
	}

	return err
}

func (s *InfoBlockService) Delete(infoBlock *models.InfoBlock) (err error) {
	filter := models.NewInfoBlockFilter()
	filter.ID = &infoBlock.ID

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][DeleteHasResourceByID] error: %v", err)
	}

	if err = s.resourceRepo.DeleteByInfoBlockID(infoBlock.ID); err != nil {
		return err
	}

	if err = s.infoBlockRepo.Delete(infoBlock); err != nil {
		return err
	}

	s.infoBlockEventService.StartJob(collection)

	return nil
}

func (s *InfoBlockService) DeleteImageFile(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if infoBlock.Image == nil {
		return infoBlock, nil
	}

	err := s.fileProvider.DeleteFile(*infoBlock.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return infoBlock, err
		}
		logger.Errorf("[info_block][InfoBlockService][DeleteImageFile] Error: %v", err)
	}
	infoBlock.Image = nil

	return s.Update(infoBlock)
}
