package service

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	app "github.com/axlle-com/blog/pkg/app/service"
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

func (s *CategoryService) SaveFromRequest(form *http.CategoryRequest, found *PostCategory, user contracts.User) (category *PostCategory, err error) {
	categoryForm := app.LoadStruct(&PostCategory{}, form).(*PostCategory)

	id := user.GetID()
	category.UserID = &id
	category.Alias = s.GenerateAlias(category)

	if found == nil {
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
		slice := make([]contracts.Gallery, 0)
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
	return s.categoryRepo.Delete(category)
}

func (s *CategoryService) Create(category *PostCategory, user contracts.User) (*PostCategory, error) {
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(category *PostCategory, found *PostCategory, user contracts.User) (*PostCategory, error) {
	if err := s.categoryRepo.Update(category, found); err != nil {
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
