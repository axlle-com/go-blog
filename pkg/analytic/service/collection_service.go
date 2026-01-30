package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/axlle-com/blog/pkg/analytic/repository"
	"github.com/google/uuid"
)

type CollectionService struct {
	api             *api.Api
	repo            repository.AnalyticRepository
	analyticService *Service
}

func NewCollectionService(
	api *api.Api,
	analyticRepo repository.AnalyticRepository,
	analyticService *Service,
) *CollectionService {
	return &CollectionService{
		api:             api,
		repo:            analyticRepo,
		analyticService: analyticService,
	}
}

func (s *CollectionService) GetAll() ([]*models.Analytic, error) {
	return s.repo.GetAll()
}

func (s *CollectionService) GetAllIds() ([]uint, error) {
	return s.repo.GetAllIds()
}

func (s *CollectionService) GetByIDs(ids []uint) ([]*models.Analytic, error) {
	return s.repo.GetByIDs(ids)
}

func (s *CollectionService) Delete(analytics []*models.Analytic) (err error) {
	var ids []uint
	for _, infoBlock := range analytics {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.repo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}

func (s *CollectionService) WithPaginate(p contract.Paginator, filter *models.AnalyticFilter) ([]*models.Analytic, error) {
	return s.repo.WithPaginate(p, filter)
}

func (s *CollectionService) Aggregates(analytics []*models.Analytic) []*models.Analytic {
	var userUUIDs []uuid.UUID

	userUUIDsMap := make(map[uuid.UUID]bool)

	for _, analytic := range analytics {
		if analytic.UserUUID != nil && *analytic.UserUUID != uuid.Nil {
			if !userUUIDsMap[*analytic.UserUUID] {
				userUUIDs = append(userUUIDs, *analytic.UserUUID)
				userUUIDsMap[*analytic.UserUUID] = true
			}
		}
	}

	var users map[uuid.UUID]contract.User

	if len(userUUIDs) > 0 {
		var err error
		users, err = s.api.User.GetMapByUUIDs(userUUIDs)
		if err != nil {
			logger.Error(err)
		}
	}

	for _, analytic := range analytics {
		if analytic.UserUUID != nil && *analytic.UserUUID != uuid.Nil {
			if user, ok := users[*analytic.UserUUID]; ok {
				analytic.User = user
			}
		}
	}

	return analytics
}
