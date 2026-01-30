package service

import (
	"errors"
	"fmt"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"gorm.io/gorm"
)

type CategoryService struct {
	api              *api.Api
	categoryRepo     repository.CategoryRepository
	aggregateService *CategoryAggregateService
}

func NewCategoryService(
	api *api.Api,
	categoryRepo repository.CategoryRepository,
) *CategoryService {
	return &CategoryService{
		api:          api,
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) SetAggregateService(aggregateService *CategoryAggregateService) {
	s.aggregateService = aggregateService
}

func (s *CategoryService) GetAggregateByID(id uint) (*models.PostCategory, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.aggregateService.Aggregate(category)
}

func (s *CategoryService) SaveFromRequest(
	form *http.CategoryRequest,
	found *models.PostCategory,
	user contract.User,
) (model *models.PostCategory, err error) {
	categoryForm := app.LoadStruct(&models.PostCategory{}, form).(*models.PostCategory)

	if found == nil {
		model, err = s.Create(categoryForm, user)
	} else {
		model, err = s.Update(categoryForm, found, user)
	}

	if err != nil {
		return
	}

	if len(form.Galleries) > 0 {
		interfaceSlice := make([]any, len(form.Galleries))
		for i, gall := range form.Galleries {
			interfaceSlice[i] = gall
		}

		slice, err := s.api.Gallery.SaveFormBatch(interfaceSlice, model)
		if err != nil {
			logger.Error(err)
		}

		model.Galleries = slice
	}

	if len(form.InfoBlocks) > 0 {
		interfaceSlice := make([]any, len(form.InfoBlocks))
		for i, block := range form.InfoBlocks {
			interfaceSlice[i] = block
		}

		slice, err := s.api.InfoBlock.CreateRelationFormBatch(interfaceSlice, model.UUID.String())
		if err != nil {
			logger.Error(err)
		}

		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *CategoryService) GetByID(id uint) (*models.PostCategory, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) Delete(category *models.PostCategory) error {
	err := s.api.Gallery.DetachResource(category)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(category)
}

func (s *CategoryService) Create(model *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	id := user.GetID()
	model.UserID = &id
	model.Alias = s.generateAlias(model)

	if err := s.categoryRepo.Create(model); err != nil {
		return nil, err
	}

	if err := s.receivedImage(model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *CategoryService) Update(model, found *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	model.ID = found.ID
	model.UUID = found.UUID
	model.UserID = found.UserID

	if model.Alias != found.Alias {
		model.Alias = s.generateAlias(model)
	}

	tx := s.categoryRepo.Tx()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.categoryRepo.WithTx(tx).Update(model, found); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if model.Image != nil && (found.Image == nil || *model.Image != *found.Image) {
		if err := s.receivedImage(model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (s *CategoryService) generateAlias(category *models.PostCategory) string {
	var aliasStr string
	if category.Alias == "" {
		aliasStr = category.Title
	} else {
		aliasStr = category.Alias
	}

	return s.api.Alias.Generate(category, aliasStr)
}

func (s *CategoryService) DeleteImageFile(category *models.PostCategory) error {
	if category.Image == nil {
		return nil
	}

	err := s.api.File.DeleteFile(*category.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	category.Image = nil

	return nil
}

func (s *CategoryService) View(category *models.PostCategory, paginator contract.Paginator, filter *models.PostFilter) (*models.PostCategory, error) {
	if s.aggregateService == nil {
		return nil, fmt.Errorf("category aggregate service is nil")
	}

	return s.aggregateService.AggregateView(category, paginator, filter)
}

// @todo переделать на фильтр везде
func (s *CategoryService) FindByParam(field string, value any) (*models.PostCategory, error) {
	return s.categoryRepo.FindByParam(field, value)
}

func (s *CategoryService) receivedImage(category *models.PostCategory) error {
	if category.Image != nil && *category.Image != "" {
		return s.api.File.Received([]string{*category.Image})
	}

	return nil
}
