package service

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
)

type PostService struct {
	postRepo          repository.PostRepository
	categoriesService *CategoriesService
	categoryService   *CategoryService
	galleryProvider   gallery.GalleryProvider
	fileProvider      provider.FileProvider
	aliasProvider     alias.AliasProvider
}

func NewPostService(
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
) *PostService {
	return &PostService{
		postRepo:          postRepo,
		categoriesService: categoriesService,
		categoryService:   categoryService,
		galleryProvider:   galleryProvider,
		fileProvider:      fileProvider,
		aliasProvider:     aliasProvider,
	}
}

func (s *PostService) GetByParam(field string, value any) (*models.Post, error) {
	return s.postRepo.GetByParam(field, value)
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	return s.postRepo.GetByID(id)
}

func (s *PostService) Update(post *models.Post) error {
	return s.postRepo.Update(post)
}
