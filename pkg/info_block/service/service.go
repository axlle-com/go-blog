package service

import (
	"errors"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	api                   *api.Api
	infoBlockRepo         repository.InfoBlockRepository
	infoBlockCollection   *CollectionService
	resourceRepo          repository.InfoBlockHasResourceRepository
	infoBlockEventService *EventService
	aggregateService      *AggregateService
}

func NewService(
	api *api.Api,
	infoBlockRepo repository.InfoBlockRepository,
	infoBlockCollection *CollectionService,
	resourceRepo repository.InfoBlockHasResourceRepository,
	infoBlockEventService *EventService,
	aggregateService *AggregateService,
) *Service {
	return &Service{
		api:                   api,
		infoBlockRepo:         infoBlockRepo,
		infoBlockCollection:   infoBlockCollection,
		resourceRepo:          resourceRepo,
		infoBlockEventService: infoBlockEventService,
		aggregateService:      aggregateService,
	}
}

func (s *Service) FindByID(id uint) (*models.InfoBlock, error) {
	return s.infoBlockRepo.FindByID(id)
}

func (s *Service) Aggregate(infoBlock *models.InfoBlock) *models.InfoBlock {
	return s.aggregateService.Aggregate(infoBlock)
}

func (s *Service) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetByIDs(ids)
}

func (s *Service) FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error) {
	return s.infoBlockRepo.FindByFilter(filter)
}

func (s *Service) SaveFromRequest(form *models.BlockRequest, found *models.InfoBlock, user contract.User) (infoBlock *models.InfoBlock, err error) {
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

func (s *Service) Create(infoBlock *models.InfoBlock, user contract.User) (*models.InfoBlock, error) {
	if user != nil {
		id := user.GetID()
		infoBlock.UserID = &id
	}

	if err := s.infoBlockRepo.Create(infoBlock); err != nil {
		return nil, err
	}

	return s.infoBlockEventService.Created(infoBlock)
}

func (s *Service) Update(new *models.InfoBlock, old *models.InfoBlock) (*models.InfoBlock, error) {
	if err := s.infoBlockRepo.Update(new, old); err != nil {
		return nil, err
	}

	return s.infoBlockEventService.Updated(new)
}

func (s *Service) Attach(resourceUUID uuid.UUID, infoBlock contract.InfoBlock) error {
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

	return s.infoBlockEventService.Attached(resourceUUID)
}

func (s *Service) DeleteByResourceUUID(resourceUUID uuid.UUID) error {
	filter := models.NewInfoBlockFilter()
	filter.ResourceUUID = &resourceUUID

	err := s.resourceRepo.DeleteByResourceUUID(resourceUUID)
	if err != nil {
		return err
	}

	return s.infoBlockEventService.DeletedByFilter(filter)
}

// @todo optimize
func (s *Service) DeleteHasResourceByID(id uint) error {
	filter := models.NewInfoBlockFilter()
	filter.RelationID = &id

	if err := s.resourceRepo.DeleteByID(id); err == nil {
		return s.infoBlockEventService.DeletedByFilter(filter)
	}

	return nil
}

func (s *Service) Delete(infoBlock *models.InfoBlock) (err error) {
	if err = s.resourceRepo.DeleteByInfoBlockID(infoBlock.ID); err != nil {
		return err
	}

	if err = s.infoBlockRepo.Delete(infoBlock); err != nil {
		return err
	}

	return s.infoBlockEventService.Deleted(infoBlock)
}

func (s *Service) DeleteImageFile(infoBlock *models.InfoBlock) (*models.InfoBlock, error) {
	if infoBlock.Image == nil {
		return infoBlock, nil
	}

	err := s.api.File.DeleteFile(*infoBlock.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return infoBlock, err
		}
		logger.Errorf("[info_block][Service][DeleteImageFile] Error: %v", err)
	}
	infoBlock.Image = nil

	oldBlock, err := s.infoBlockRepo.FindByID(infoBlock.ID)
	if err != nil {
		return infoBlock, err
	}

	return s.Update(infoBlock, oldBlock)
}
