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

type PostAggregateService struct {
	api                  *api.Api
	postRepo             repository.PostRepository
	categoriesService    *CategoryCollectionService
	categoryService      *CategoryService
	tagCollectionService *TagCollectionService
}

func NewPostAggregateService(
	api *api.Api,
	postRepo repository.PostRepository,
	categoriesService *CategoryCollectionService,
	categoryService *CategoryService,
	tagCollectionService *TagCollectionService,
) *PostAggregateService {
	return &PostAggregateService{
		api:                  api,
		postRepo:             postRepo,
		categoriesService:    categoriesService,
		categoryService:      categoryService,
		tagCollectionService: tagCollectionService,
	}
}

func (s *PostAggregateService) Aggregate(post *models.Post) (*models.Post, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)
	var tags = make([]*models.PostTag, 0)
	var err error

	service.SafeGo(&wg, func() {
		galleries = s.api.Gallery.GetForResourceUUID(post.UUID.String())
	})

	service.SafeGo(&wg, func() {
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(post.UUID.String())
	})

	service.SafeGo(&wg, func() {
		tags, err = s.tagCollectionService.GetForResource(post)
	})

	wg.Wait()

	post.Galleries = galleries
	post.InfoBlocks = infoBlocks
	post.PostTags = tags

	return post, err
}

func (s *PostAggregateService) AggregateView(post *models.Post) (*models.Post, error) {
	var wg sync.WaitGroup
	agg := errutil.New()

	s.addInfoBlocks(post, &wg, agg)
	s.addGalleries(post, &wg, agg)
	s.addTags(post, &wg, agg)

	wg.Wait()

	return post, agg.ErrorAndReset()
}

func (s *PostAggregateService) addInfoBlocks(post *models.Post, wg *sync.WaitGroup, agg *errutil.ErrUtil) {
	if post.InfoBlocksSnapshot == nil {
		service.SafeGo(wg, func(p *models.Post, id uuid.UUID) func() {
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

	if len(post.InfoBlocksSnapshot) > 2 {
		var blocks []dto.InfoBlock
		if err := json.Unmarshal(post.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][blogController][RenderPost] id=%v: %v", post.ID, err)
		} else {
			var interfaceBlocks []contract.InfoBlock
			for _, block := range blocks {
				interfaceBlocks = append(interfaceBlocks, block)
			}

			post.InfoBlocks = interfaceBlocks
		}
	}
}

func (s *PostAggregateService) addGalleries(post *models.Post, wg *sync.WaitGroup, agg *errutil.ErrUtil) {
	if post.GalleriesSnapshot == nil {
		service.SafeGo(wg, func(p *models.Post, id uuid.UUID) func() {
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

	if len(post.GalleriesSnapshot) > 2 {
		var galleries []dto.Gallery
		if err := json.Unmarshal(post.GalleriesSnapshot, &galleries); err != nil {
			logger.Errorf("[blog][blogController][RenderPost] id=%v: %v", post.ID, err)
		} else {
			interfaceGalleries := make([]contract.Gallery, 0, len(galleries))
			for _, gallery := range galleries {
				interfaceGalleries = append(interfaceGalleries, gallery)
			}

			post.Galleries = interfaceGalleries
		}
	}
}

func (s *PostAggregateService) addTags(post *models.Post, wg *sync.WaitGroup, agg *errutil.ErrUtil) {
	service.SafeGo(wg, func(p *models.Post) func() {
		return func() {
			ts, e := s.tagCollectionService.GetForResource(p)
			if e != nil {
				agg.Add(fmt.Errorf("get tags: %w", e))
				return
			}

			p.PostTags = ts
		}
	}(post))
}
