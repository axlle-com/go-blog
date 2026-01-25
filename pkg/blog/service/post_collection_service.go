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

type PostCollectionService struct {
	postRepo          repository.PostRepository
	categoriesService *CategoriesService
	categoryService   *CategoryService
	api               *api.Api
}

func NewPostCollectionService(
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	api *api.Api,
) *PostCollectionService {
	return &PostCollectionService{
		postRepo:          postRepo,
		categoriesService: categoriesService,
		categoryService:   categoryService,
		api:               api,
	}
}

func (s *PostCollectionService) Aggregates(posts []*models.Post) []*models.Post {
	var userIDs []uint
	var categoryIDs []uint
	var templateNames []string

	userIDsMap := make(map[uint]bool)
	categoryIDsMap := make(map[uint]bool)
	templateNamesMap := make(map[string]bool)

	for _, post := range posts {
		if post.TemplateName != "" && !templateNamesMap[post.TemplateName] {
			templateNames = append(templateNames, post.TemplateName)
			templateNamesMap[post.TemplateName] = true
		}
		if post.UserID != nil {
			id := *post.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
		if post.PostCategoryID != nil {
			id := *post.PostCategoryID
			if !categoryIDsMap[id] {
				categoryIDs = append(categoryIDs, id)
				categoryIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User
	var templates map[string]contract.Template
	var categories map[uint]*models.PostCategory

	wg.Add(3)

	go func() {
		defer wg.Done()
		if len(templateNames) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByNames(templateNames)
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
			categories, err = s.categoriesService.GetMapByIDs(categoryIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, post := range posts {
		if templates != nil && post.TemplateName != "" {
			post.Template = templates[post.TemplateName]
		}
		if post.UserID != nil {
			post.User = users[*post.UserID]
		}
		if post.PostCategoryID != nil {
			post.Category = categories[*post.PostCategoryID]
		}
	}

	return posts
}

func (s *PostCollectionService) WithPaginate(p contract.Paginator, filter *models.PostFilter) ([]*models.Post, error) {
	return s.postRepo.WithPaginate(p, filter)
}

func (s *PostCollectionService) GetAll() ([]*models.Post, error) {
	return s.postRepo.GetAll()
}

func (s *PostCollectionService) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	return s.postRepo.UpdateFieldsByUUIDs(uuids, patch)
}
