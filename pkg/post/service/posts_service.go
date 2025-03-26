package service

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
)

type PostsService struct {
	postRepo          repository.PostRepository
	categoriesService *CategoriesService
	categoryService   *CategoryService
	galleryProvider   gallery.GalleryProvider
	fileProvider      provider.FileProvider
	aliasProvider     alias.AliasProvider
}

func NewPostsService(
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
) *PostsService {
	return &PostsService{
		postRepo:          postRepo,
		categoriesService: categoriesService,
		categoryService:   categoryService,
		galleryProvider:   galleryProvider,
		fileProvider:      fileProvider,
		aliasProvider:     aliasProvider,
	}
}

func (s *PostsService) WithPaginate(p contracts.Paginator, filter *models.PostFilter) ([]*models.PostFull, error) {
	return s.postRepo.WithPaginate(p, filter)
}
