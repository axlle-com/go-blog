package service

import (
	"errors"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/alias"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/models"
	. "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	infoBlockProvider "github.com/axlle-com/blog/pkg/info_block/provider"
	"gorm.io/gorm"
	"sync"
)

type CategoryService struct {
	categoryRepo      repository.CategoryRepository
	galleryProvider   gallery.GalleryProvider
	fileProvider      provider.FileProvider
	aliasProvider     alias.AliasProvider
	infoBlockProvider infoBlockProvider.InfoBlockProvider
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	aliasProvider alias.AliasProvider,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	infoBlockProvider infoBlockProvider.InfoBlockProvider,
) *CategoryService {
	return &CategoryService{
		categoryRepo:      categoryRepo,
		galleryProvider:   galleryProvider,
		fileProvider:      fileProvider,
		aliasProvider:     aliasProvider,
		infoBlockProvider: infoBlockProvider,
	}
}

func (s *CategoryService) GetAggregateByID(id uint) (*PostCategory, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.Aggregate(category)
}

func (s *CategoryService) Aggregate(category *PostCategory) (*PostCategory, error) {
	var wg sync.WaitGroup

	var galleries = make([]contracts.Gallery, 0)
	var infoBlocks = make([]contracts.InfoBlock, 0)

	wg.Add(2)

	go func() {
		defer wg.Done()
		galleries = s.galleryProvider.GetForResource(category)
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.infoBlockProvider.GetForResource(category)
	}()

	wg.Wait()

	category.Galleries = galleries
	category.InfoBlocks = infoBlocks

	return category, nil
}

func (s *CategoryService) SaveFromRequest(form *http.CategoryRequest, found *PostCategory, user contracts.User) (model *PostCategory, err error) {
	categoryForm := app.LoadStruct(&PostCategory{}, form).(*PostCategory)

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

		slice, err := s.infoBlockProvider.SaveFormBatch(interfaceSlice, model)
		if err != nil {
			logger.Error(err)
		}
		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *CategoryService) GetByID(id uint) (*PostCategory, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) Delete(category *PostCategory) error {
	err := s.galleryProvider.DetachResource(category)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(category)
}

func (s *CategoryService) Create(category *PostCategory, user contracts.User) (*PostCategory, error) {
	id := user.GetID()
	category.UserID = &id
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(category *PostCategory, found *PostCategory, user contracts.User) (*PostCategory, error) {
	category.ID = found.ID
	category.UUID = found.UUID

	tx := db.GetDB().Begin()

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

func (s *CategoryService) GenerateAlias(category *PostCategory) string {
	var aliasStr string
	if category.Alias == "" {
		aliasStr = category.Title
	} else {
		aliasStr = category.Alias
	}

	return s.aliasProvider.Generate(category, aliasStr)
}

func (s *CategoryService) DeleteImageFile(category *PostCategory) error {
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
