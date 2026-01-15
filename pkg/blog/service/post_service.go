package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PostService struct {
	queue                contract.Queue
	postRepo             repository.PostRepository
	categoriesService    *CategoriesService
	categoryService      *CategoryService
	tagCollectionService *TagCollectionService
	api                  *api.Api
}

func NewPostService(
	queue contract.Queue,
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	tagCollectionService *TagCollectionService,
	api *api.Api,
) *PostService {
	return &PostService{
		queue:                queue,
		postRepo:             postRepo,
		categoriesService:    categoriesService,
		categoryService:      categoryService,
		tagCollectionService: tagCollectionService,
		api:                  api,
	}
}

func (s *PostService) FindAggregateByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
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
		galleries = s.api.Gallery.GetForResourceUUID(post.UUID.String())
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(post.UUID.String())
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

func (s *PostService) View(post *models.Post) (*models.Post, error) {
	var wg sync.WaitGroup
	agg := errutil.New()

	// --- InfoBlocks snapshot ---
	if post.InfoBlocksSnapshot == nil {
		service.SafeGo(&wg, func(p *models.Post, id uuid.UUID) func() {
			return func() {
				blocks := s.api.InfoBlock.GetForResourceUUID(id.String())

				mapped := dto.MapInfoBlocks(blocks)
				if mapped == nil {
					mapped = []dto.InfoBlock{}
				}

				raw, e := json.Marshal(mapped)
				if e != nil {
					agg.Add(fmt.Errorf("marshal info_blocks_snapshot: %w", e))
					return
				}

				v := datatypes.JSON(raw)
				patch := map[string]any{"info_blocks_snapshot": v}
				if _, e = s.postRepo.UpdateFieldsByUUIDs([]uuid.UUID{id}, patch); e != nil {
					agg.Add(fmt.Errorf("update info_blocks_snapshot: %w", e))
					return
				}

				p.InfoBlocksSnapshot = v
			}
		}(post, post.UUID))
	}

	// --- Galleries snapshot ---
	if post.GalleriesSnapshot == nil {
		service.SafeGo(&wg, func(p *models.Post, id uuid.UUID) func() {
			return func() {
				galleries := s.api.Gallery.GetForResourceUUID(id.String())

				mapped := dto.MapGalleries(galleries)
				if mapped == nil {
					mapped = []dto.Gallery{}
				}

				raw, e := json.Marshal(mapped)
				if e != nil {
					agg.Add(fmt.Errorf("marshal galleries_snapshot: %w", e))
					return
				}

				v := datatypes.JSON(raw)
				patch := map[string]any{"galleries_snapshot": v}
				if _, e = s.postRepo.UpdateFieldsByUUIDs([]uuid.UUID{id}, patch); e != nil {
					agg.Add(fmt.Errorf("update galleries_snapshot: %w", e))
					return
				}

				p.GalleriesSnapshot = v
			}
		}(post, post.UUID))
	}

	// --- Tags ---
	service.SafeGo(&wg, func(p *models.Post) func() {
		return func() {
			ts, e := s.tagCollectionService.GetForResource(p)
			if e != nil {
				agg.Add(fmt.Errorf("get tags: %w", e))
				return
			}
			p.PostTags = ts
		}
	}(post))

	// --- template ---
	service.SafeGo(&wg, func(p *models.Post) func() {
		return func() {
			tpl, e := s.api.Template.GetByID(*p.TemplateID)
			if e != nil {
				agg.Add(fmt.Errorf("get template: %w", e))
				return
			}
			p.Template = tpl
		}
	}(post))

	wg.Wait()

	return post, agg.Error()
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
