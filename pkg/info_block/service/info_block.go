package service

import (
	"errors"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InfoBlockService struct {
	infoBlockRepo         repository.InfoBlockRepository
	infoBlockCollection   *InfoBlockCollectionService
	resourceRepo          repository.InfoBlockHasResourceRepository
	infoBlockEventService *InfoBlockEventService
	api                   *api.Api
}

func NewInfoBlockService(
	infoBlockRepo repository.InfoBlockRepository,
	infoBlockCollection *InfoBlockCollectionService,
	resourceRepo repository.InfoBlockHasResourceRepository,
	infoBlockEventService *InfoBlockEventService,
	api *api.Api,
) *InfoBlockService {
	return &InfoBlockService{
		infoBlockRepo:         infoBlockRepo,
		infoBlockCollection:   infoBlockCollection,
		resourceRepo:          resourceRepo,
		infoBlockEventService: infoBlockEventService,
		api:                   api,
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
			infoBlock.Template, err = s.api.Template.GetByID(*infoBlock.TemplateID)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if infoBlock.UserID != nil && *infoBlock.UserID != 0 {
			var err error
			infoBlock.User, err = s.api.User.GetByID(*infoBlock.UserID)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		infoBlock.Galleries = s.api.Gallery.GetForResourceUUID(infoBlock.UUID.String())
	}()

	wg.Wait()

	return infoBlock
}

func (s *InfoBlockService) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetByIDs(ids)
}

func (s *InfoBlockService) FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error) {
	return s.infoBlockRepo.FindByFilter(filter)
}

func (s *InfoBlockService) SaveFromRequest(form *models.BlockRequest, found *models.InfoBlock, user contract.User) (infoBlock *models.InfoBlock, err error) {
	blockForm := app.LoadStruct(&models.InfoBlock{}, form).(*models.InfoBlock)

	if found == nil {
		infoBlock, err = s.Create(blockForm, user)
	} else {
		blockForm.ID = found.ID
		blockForm.UUID = found.UUID
		infoBlock, err = s.Update(blockForm, found)
	}

	if err != nil {
		return
	}

	if len(form.Galleries) > 0 {
		slice := make([]contract.Gallery, 0)
		for _, gRequest := range form.Galleries {
			if gRequest == nil {
				continue
			}

			g, err := s.api.Gallery.SaveForm(gRequest, infoBlock)
			if err != nil || g == nil {
				continue
			}
			slice = append(slice, g)
		}
		infoBlock.Galleries = slice
	}
	return
}

func (s *InfoBlockService) Create(infoBlock *models.InfoBlock, user contract.User) (*models.InfoBlock, error) {
	if user != nil {
		id := user.GetID()
		infoBlock.UserID = &id
	}
	if err := s.infoBlockRepo.Create(infoBlock); err != nil {
		return nil, err
	}

	return infoBlock, nil
}

func (s *InfoBlockService) Update(new *models.InfoBlock, old *models.InfoBlock) (*models.InfoBlock, error) {
	if err := s.infoBlockRepo.Update(new, old); err != nil {
		return nil, err
	}

	filter := models.NewInfoBlockFilter()
	filter.ID = &new.ID

	collection, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockService][Update] error: %v", err)
	}

	s.infoBlockEventService.StartJob(collection, "update")

	return new, nil
}

func (s *InfoBlockService) Attach(resourceUUID uuid.UUID, infoBlock contract.InfoBlock) error {
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

	s.infoBlockEventService.StartJob(collection, "attach")

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

	s.infoBlockEventService.StartJob(collection, "delete")

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
		s.infoBlockEventService.StartJob(collection, "delete")
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

	s.infoBlockEventService.StartJob(collection, "delete")

	return nil
}

func (s *InfoBlockService) DeleteImageFile(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if infoBlock.Image == nil {
		return infoBlock, nil
	}

	err := s.api.File.DeleteFile(*infoBlock.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return infoBlock, err
		}
		logger.Errorf("[info_block][InfoBlockService][DeleteImageFile] Error: %v", err)
	}
	infoBlock.Image = nil

	oldBlock, err := s.infoBlockRepo.FindByID(infoBlock.ID)
	if err != nil {
		return infoBlock, err
	}

	return s.Update(infoBlock, oldBlock)
}
