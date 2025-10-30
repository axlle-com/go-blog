package service

import (
	"sync"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/file/provider"
)

type PostService struct {
	queue                contract.Queue
	postRepo             repository.PostRepository
	categoriesService    *CategoriesService
	categoryService      *CategoryService
	tagCollectionService *TagCollectionService
	galleryProvider      appPovider.GalleryProvider
	fileProvider         provider.FileProvider
	aliasProvider        alias.AliasProvider
	infoBlockProvider    appPovider.InfoBlockProvider
}

func NewPostService(
	queue contract.Queue,
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	tagCollectionService *TagCollectionService,
	galleryProvider appPovider.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
	infoBlockProvider appPovider.InfoBlockProvider,
) *PostService {
	return &PostService{
		queue:                queue,
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

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)
	var tags = make([]*models.PostTag, 0)
	var err error

	wg.Add(3)

	go func() {
		defer wg.Done()
		galleries = s.galleryProvider.GetForResourceUUID(post.UUID.String())
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.infoBlockProvider.GetForResourceUUID(post.UUID.String())
	}()

	go func() {
		defer wg.Done()
		tags, err = s.tagCollectionService.GetForResource(post)
		if err != nil {
			logger.Errorf("[PostService] Error: %v", err)
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

func (s *PostService) generateAlias(post *models.Post) string {
	var alias string
	if post.Alias == "" {
		alias = post.Title
	} else {
		alias = post.Alias
	}

	return s.aliasProvider.Generate(post, alias)
}

func (s *PostService) receivedImage(post *models.Post) error {
	if post.Image != nil && *post.Image != "" {
		return s.fileProvider.Received([]string{*post.Image})
	}

	return nil
}
