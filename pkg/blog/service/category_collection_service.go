package service

import (
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/google/uuid"
)

type CategoryCollectionService struct {
	categoryRepo repository.CategoryRepository
	api          *api.Api
}

func NewCategoryCollectionService(
	categoryRepo repository.CategoryRepository,
	api *api.Api,
) *CategoryCollectionService {
	return &CategoryCollectionService{
		categoryRepo: categoryRepo,
		api:          api,
	}
}

func (s *CategoryCollectionService) GetAggregates(categories []*models.PostCategory) []*models.PostCategory {
	var templateNames []string
	var userIDs []uint
	var categoryIDs []uint

	templateNamesMap := make(map[string]bool)
	userIDsMap := make(map[uint]bool)
	categoryIDsMap := make(map[uint]bool)

	for _, category := range categories {
		if category.TemplateName != "" && !templateNamesMap[category.TemplateName] {
			templateNames = append(templateNames, category.TemplateName)
			templateNamesMap[category.TemplateName] = true
		}

		if category.UserID != nil {
			id := *category.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}

		if category.PostCategoryID != nil {
			id := *category.PostCategoryID
			if !categoryIDsMap[id] {
				categoryIDs = append(categoryIDs, id)
				categoryIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User
	var templates map[string]contract.Template
	var parents map[uint]*models.PostCategory

	service.SafeGo(&wg, func() {
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	})

	service.SafeGo(&wg, func() {
		if len(categoryIDs) > 0 {
			var err error
			parents, err = s.GetMapByIDs(categoryIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	})

	service.SafeGo(&wg, func() {
		if len(templateNames) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByNames(templateNames)
			if err != nil {
				logger.Error(err)
			}
		}
	})

	wg.Wait()

	for _, category := range categories {
		if templates != nil && category.TemplateName != "" {
			category.Template = templates[category.TemplateName]
		}

		if category.UserID != nil {
			category.User = users[*category.UserID]
		}

		if category.PostCategoryID != nil {
			category.Category = parents[*category.PostCategoryID]
		}
	}

	return categories
}

func (s *CategoryCollectionService) GetAll() ([]*models.PostCategory, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoryCollectionService) GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error) {
	return s.categoryRepo.GetAllForParent(parent)
}

func (s *CategoryCollectionService) WithPaginate(p contract.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error) {
	return s.categoryRepo.WithPaginate(p, filter)
}

func (s *CategoryCollectionService) GetMapByIDs(ids []uint) (map[uint]*models.PostCategory, error) {
	categories, err := s.categoryRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
	collection := make(map[uint]*models.PostCategory, len(categories))
	for _, item := range categories {
		collection[item.ID] = item
	}
	return collection, nil
}

func (s *CategoryCollectionService) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	return s.categoryRepo.UpdateFieldsByUUIDs(uuids, patch)
}
