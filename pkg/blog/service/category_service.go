package service

import (
	"errors"
	"sync"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/alias"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/file/provider"
	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo      repository.CategoryRepository
	galleryProvider   appPovider.GalleryProvider
	fileProvider      provider.FileProvider
	aliasProvider     alias.AliasProvider
	infoBlockProvider appPovider.InfoBlockProvider
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	aliasProvider alias.AliasProvider,
	galleryProvider appPovider.GalleryProvider,
	fileProvider provider.FileProvider,
	infoBlockProvider appPovider.InfoBlockProvider,
) *CategoryService {
	return &CategoryService{
		categoryRepo:      categoryRepo,
		galleryProvider:   galleryProvider,
		fileProvider:      fileProvider,
		aliasProvider:     aliasProvider,
		infoBlockProvider: infoBlockProvider,
	}
}

func (s *CategoryService) GetAggregateByID(id uint) (*models.PostCategory, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.Aggregate(category)
}

func (s *CategoryService) Aggregate(category *models.PostCategory) (*models.PostCategory, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)

	wg.Add(2)

	go func() {
		defer wg.Done()
		galleries = s.galleryProvider.GetForResourceUUID(category.UUID.String())
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.infoBlockProvider.GetForResourceUUID(category.UUID.String())
	}()

	wg.Wait()

	category.Galleries = galleries
	category.InfoBlocks = infoBlocks

	return category, nil
}

func (s *CategoryService) SaveFromRequest(
	form *http.CategoryRequest,
	found *models.PostCategory,
	user contract.User,
) (model *models.PostCategory, err error) {
	categoryForm := app.LoadStruct(&models.PostCategory{}, form).(*models.PostCategory)

	categoryForm.Alias = s.GenerateAlias(categoryForm)

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

		slice, err := s.galleryProvider.SaveFormBatch(interfaceSlice, model)
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

		slice, err := s.infoBlockProvider.SaveFormBatch(interfaceSlice, model.UUID.String())
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
	err := s.galleryProvider.DetachResource(category)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(category)
}

func (s *CategoryService) Create(category *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	id := user.GetID()
	category.UserID = &id
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(category *models.PostCategory, found *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	tx := s.categoryRepo.Tx()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.categoryRepo.WithTx(tx).Update(category, found); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GenerateAlias(category *models.PostCategory) string {
	var aliasStr string
	if category.Alias == "" {
		aliasStr = category.Title
	} else {
		aliasStr = category.Alias
	}

	return s.aliasProvider.Generate(category, aliasStr)
}

func (s *CategoryService) DeleteImageFile(category *models.PostCategory) error {
	if category.Image == nil {
		return nil
	}
	err := s.fileProvider.DeleteFile(*category.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		logger.Errorf("[DeleteImageFile] Error: %v", err)
	}
	category.Image = nil
	return nil
}
