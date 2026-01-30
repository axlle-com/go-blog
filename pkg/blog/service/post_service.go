package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/queue/job"
	"github.com/axlle-com/blog/pkg/blog/repository"
)

type PostService struct {
	queue contract.Queue
	api   *api.Api

	postRepo             repository.PostRepository
	postAggregateService *PostAggregateService
	categoriesService    *CategoryCollectionService
	categoryService      *CategoryService
	tagCollectionService *TagCollectionService
}

func NewPostService(
	queue contract.Queue,
	api *api.Api,
	postRepo repository.PostRepository,
	postAggregateService *PostAggregateService,
	categoriesService *CategoryCollectionService,
	categoryService *CategoryService,
	tagCollectionService *TagCollectionService,
) *PostService {
	return &PostService{
		queue:                queue,
		api:                  api,
		postRepo:             postRepo,
		postAggregateService: postAggregateService,
		categoriesService:    categoriesService,
		categoryService:      categoryService,
		tagCollectionService: tagCollectionService,
	}
}

func (s *PostService) FindAggregateByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.postAggregateService.Aggregate(post)
}

func (s *PostService) View(post *models.Post) (*models.Post, error) {
	return s.postAggregateService.AggregateView(post)
}

func (s *PostService) FindByParam(field string, value any) (*models.Post, error) {
	return s.postRepo.FindByParam(field, value)
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	return s.postRepo.FindByID(id)
}

func (s *PostService) generateAlias(model *models.Post) string {
	var newAlias string

	if model.IsMain {
		return newAlias
	}

	if model.Alias == "" {
		newAlias = model.Title
	} else {
		newAlias = model.Alias
	}

	return s.api.Alias.Generate(model, newAlias)
}

func (s *PostService) receivedImage(model *models.Post) error {
	if model.Image != nil && *model.Image != "" {
		return s.api.File.Received([]string{*model.Image})
	}

	return nil
}

func (s *PostService) PostDelete(post *models.Post) error {
	err := s.api.Gallery.DetachResource(post)
	if err != nil {
		return err
	}

	err = s.api.InfoBlock.DetachResourceUUID(post.UUID.String())
	if err != nil {
		return err
	}

	if err := s.postRepo.Delete(post); err != nil {
		return err
	}

	s.queue.Enqueue(job.NewPostJob(post, "delete"), 0)

	return nil
}
