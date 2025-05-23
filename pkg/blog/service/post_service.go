package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	provider2 "github.com/axlle-com/blog/pkg/info_block/provider"
	"sync"
)

type PostService struct {
	postRepo             repository.PostRepository
	categoriesService    *CategoriesService
	categoryService      *CategoryService
	tagCollectionService *TagCollectionService
	galleryProvider      gallery.GalleryProvider
	fileProvider         provider.FileProvider
	aliasProvider        alias.AliasProvider
	infoBlockProvider    provider2.InfoBlockProvider
}

func NewPostService(
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	tagCollectionService *TagCollectionService,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
	infoBlockProvider provider2.InfoBlockProvider,
) *PostService {
	return &PostService{
		postRepo:             postRepo,
		categoriesService:    categoriesService,
		categoryService:      categoryService,
		tagCollectionService: tagCollectionService,
		galleryProvider:      galleryProvider,
		fileProvider:         fileProvider,
		aliasProvider:        aliasProvider,
		infoBlockProvider:    infoBlockProvider,
	}
}

func (s *PostService) GetAggregateByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.Aggregate(post)
}

func (s *PostService) Aggregate(post *models.Post) (*models.Post, error) {
	var wg sync.WaitGroup

	var galleries = make([]contracts.Gallery, 0)
	var infoBlocks = make([]contracts.InfoBlock, 0)
	var tags = make([]*models.PostTag, 0)
	var err error

	wg.Add(3)

	go func() {
		defer wg.Done()
		galleries = s.galleryProvider.GetForResource(post)
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.infoBlockProvider.GetForResource(post)
	}()

	go func() {
		defer wg.Done()
		tags, err = s.tagCollectionService.GetForResource(post)
		if err != nil {
			logger.Errorf("[] Error: %v", err)
		}
	}()

	wg.Wait()

	post.Galleries = galleries
	post.InfoBlocks = infoBlocks
	post.PostTags = tags

	return post, nil
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
