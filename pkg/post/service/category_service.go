package service

import (
	"errors"
	"github.com/axlle-com/blog/app/db"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
)

type CategoryService struct {
	categoryRepo    repository.CategoryRepository
	galleryProvider gallery.GalleryProvider
	fileProvider    provider.FileProvider
	aliasProvider   alias.AliasProvider
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	aliasProvider alias.AliasProvider,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
) *CategoryService {
	return &CategoryService{
		categoryRepo:    categoryRepo,
		galleryProvider: galleryProvider,
		fileProvider:    fileProvider,
		aliasProvider:   aliasProvider,
	}
}

func (s *CategoryService) SaveFromRequest(form *http.CategoryRequest, found *PostCategory, user contracts2.User) (category *PostCategory, err error) {
	categoryForm := app.LoadStruct(&PostCategory{}, form).(*PostCategory)

	categoryForm.Alias = s.GenerateAlias(categoryForm)

	if found == nil {
		id := user.GetID()
		categoryForm.UserID = &id
		category, err = s.Create(categoryForm, user)
	} else {
		categoryForm.ID = found.ID
		categoryForm.UUID = found.UUID
		category, err = s.Update(categoryForm, found, user)
	}

	if err != nil {
		return
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts2.Gallery, 0)
		for _, gRequest := range form.Galleries {
			if gRequest == nil {
				continue
			}

			g, err := s.galleryProvider.SaveFromForm(gRequest, category)
			if err != nil || g == nil {
				continue
			}
			slice = append(slice, g)
		}
		category.Galleries = slice
	}
	return
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

func (s *CategoryService) Create(category *PostCategory, user contracts2.User) (*PostCategory, error) {
	id := user.GetID()
	category.UserID = &id
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(category *PostCategory, found *PostCategory, user contracts2.User) (*PostCategory, error) {
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
		return errors.New("image is nil")
	}
	err := s.fileProvider.DeleteFile(*category.Image)
	if err != nil {
		return err
	}
	category.Image = nil
	return nil
}
