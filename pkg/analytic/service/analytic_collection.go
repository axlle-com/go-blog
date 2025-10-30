package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/axlle-com/blog/pkg/analytic/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
)

type AnalyticCollectionService struct {
	analyticRepo    repository.AnalyticRepository
	analyticService *AnalyticService
	userProvider    userProvider.UserProvider
}

func NewAnalyticCollectionService(
	analyticRepo repository.AnalyticRepository,
	analyticService *AnalyticService,
	userProvider userProvider.UserProvider,
) *AnalyticCollectionService {
	return &AnalyticCollectionService{
		analyticService: analyticService,
		analyticRepo:    analyticRepo,
		userProvider:    userProvider,
	}
}

func (s *AnalyticCollectionService) GetAll() ([]*models.Analytic, error) {
	return s.analyticRepo.GetAll()
}

func (s *AnalyticCollectionService) GetAllIds() ([]uint, error) {
	return s.analyticRepo.GetAllIds()
}

func (s *AnalyticCollectionService) GetByIDs(ids []uint) ([]*models.Analytic, error) {
	return s.analyticRepo.GetByIDs(ids)
}

func (s *AnalyticCollectionService) Delete(analytics []*models.Analytic) (err error) {
	var ids []uint
	for _, infoBlock := range analytics {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.analyticRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}

func (s *AnalyticCollectionService) WithPaginate(p contract.Paginator, filter *models.AnalyticFilter) ([]*models.Analytic, error) {
	return s.analyticRepo.WithPaginate(p, filter)
}

func (s *AnalyticCollectionService) Aggregates(analytics []*models.Analytic) []*models.Analytic {
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
		users, err = s.userProvider.GetMapByUUIDs(userUUIDs)
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
