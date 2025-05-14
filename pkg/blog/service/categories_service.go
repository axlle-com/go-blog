package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/alias"
	. "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"sync"
)

type CategoriesService struct {
	categoryRepo    repository.CategoryRepository
	template        template.TemplateProvider
	user            user.UserProvider
	galleryProvider gallery.GalleryProvider
	aliasProvider   alias.AliasProvider
}

func NewCategoriesService(
	categoryRepo repository.CategoryRepository,
	aliasProvider alias.AliasProvider,
	galleryProvider gallery.GalleryProvider,
	template template.TemplateProvider,
	user user.UserProvider,
) *CategoriesService {
	return &CategoriesService{
		categoryRepo:    categoryRepo,
		template:        template,
		user:            user,
		galleryProvider: galleryProvider,
		aliasProvider:   aliasProvider,
	}
}

func (s *CategoriesService) GetAggregates(categories []*PostCategory) []*PostCategory {
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

	var users map[uint]contracts.User
	var templates map[uint]contracts.Template
	var parents map[uint]*PostCategory

	wg.Add(3)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.template.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.user.GetMapByIDs(userIDs)
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

func (s *CategoriesService) GetAll() ([]*PostCategory, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoriesService) GetAllForParent(parent *PostCategory) ([]*PostCategory, error) {
	return s.categoryRepo.GetAllForParent(parent)
}

func (s *CategoriesService) WithPaginate(p contracts.Paginator, filter *CategoryFilter) ([]*PostCategory, error) {
	return s.categoryRepo.WithPaginate(p, filter)
}

func (s *CategoriesService) GetMapByIDs(ids []uint) (map[uint]*PostCategory, error) {
	categories, err := s.categoryRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
	collection := make(map[uint]*PostCategory, len(categories))
	for _, item := range categories {
		collection[item.ID] = item
	}
	return collection, nil
}
