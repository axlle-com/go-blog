package service

import (
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/google/uuid"
)

type CategoriesService struct {
	categoryRepo repository.CategoryRepository
	api          *api.Api
}

func NewCategoriesService(
	categoryRepo repository.CategoryRepository,
	api *api.Api,
) *CategoriesService {
	return &CategoriesService{
		categoryRepo: categoryRepo,
		api:          api,
	}
}

func (s *CategoriesService) GetAggregates(categories []*models.PostCategory) []*models.PostCategory {
	var templateIDs []uint
	var userIDs []uint
	var categoryIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
	categoryIDsMap := make(map[uint]bool)

	for _, category := range categories {
		if category.TemplateID != nil {
			id := *category.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
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
	var templates map[uint]contract.Template
	var parents map[uint]*models.PostCategory

	wg.Add(3)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(categoryIDs) > 0 {
			var err error
			parents, err = s.GetMapByIDs(categoryIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, category := range categories {
		if category.TemplateID != nil {
			category.Template = templates[*category.TemplateID]
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

func (s *CategoriesService) GetAll() ([]*models.PostCategory, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoriesService) GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error) {
	return s.categoryRepo.GetAllForParent(parent)
}

func (s *CategoriesService) WithPaginate(p contract.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error) {
	return s.categoryRepo.WithPaginate(p, filter)
}

func (s *CategoriesService) GetMapByIDs(ids []uint) (map[uint]*models.PostCategory, error) {
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

func (s *CategoriesService) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	return s.categoryRepo.UpdateFieldsByUUIDs(uuids, patch)
}
