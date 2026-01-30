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

type CategoryAggregateService struct {
	api                   *api.Api
	categoryRepo          repository.CategoryRepository
	postCollectionService *PostCollectionService
}

func NewCategoryAggregateService(
	api *api.Api,
	categoryRepo repository.CategoryRepository,
	postCollectionService *PostCollectionService,
) *CategoryAggregateService {
	return &CategoryAggregateService{
		api:                   api,
		categoryRepo:          categoryRepo,
		postCollectionService: postCollectionService,
	}
}

func (s *CategoryAggregateService) Aggregate(category *models.PostCategory) (*models.PostCategory, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)

	service.SafeGo(&wg, func() {
		galleries = s.api.Gallery.GetForResourceUUID(category.UUID.String())
	})

	service.SafeGo(&wg, func() {
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(category.UUID.String())
	})

	service.SafeGo(&wg, func() {
		if category.TemplateName == "" {
			return
		}

		tpl, e := s.api.Template.GetByName(category.TemplateName)
		if e != nil {
			logger.Errorf("[CategoryAggregateService] Error: %v", e)
			return
		}

		category.Template = tpl
	})

	wg.Wait()

	category.Galleries = galleries
	category.InfoBlocks = infoBlocks

	return category, nil
}

func (s *CategoryAggregateService) AggregateView(
	category *models.PostCategory,
	paginator contract.Paginator,
	filter *models.PostFilter,
) (*models.PostCategory, error) {
	var wg sync.WaitGroup
	agg := errutil.New()

	s.addPosts(category, paginator, filter, &wg, agg)
	s.addInfoBlocks(category, &wg, agg)
	s.addGalleries(category, &wg, agg)

	wg.Wait()

	return category, agg.Error()
}

func (s *CategoryAggregateService) addPosts(
	category *models.PostCategory,
	paginator contract.Paginator,
	filter *models.PostFilter,
	wg *sync.WaitGroup,
	agg *errutil.ErrUtil,
) {
	if s.postCollectionService != nil && paginator != nil {
		service.SafeGo(wg, func(c *models.PostCategory, p contract.Paginator, f *models.PostFilter) func() {
			return func() {
				postFilter := models.NewPostFilter()
				if f != nil {
					*postFilter = *f
				}
				categoryID := c.ID
				postFilter.PostCategoryID = &categoryID

				posts, e := s.postCollectionService.WithPaginate(p, postFilter)
				if e != nil {
					agg.Add(fmt.Errorf("get posts: %w", e))
					return
				}

				c.Posts = posts
			}
		}(category, paginator, filter))
	}
}

func (s *CategoryAggregateService) addInfoBlocks(category *models.PostCategory, wg *sync.WaitGroup, agg *errutil.ErrUtil) {
	if category.InfoBlocksSnapshot == nil {
		service.SafeGo(wg, func(c *models.PostCategory, id uuid.UUID) func() {
			return func() {
				blocks := s.api.InfoBlock.GetForResourceUUID(id.String())
				if len(blocks) == 0 {
					return
				}

				raw, e := json.Marshal(dto.MapInfoBlocks(blocks))
				if e != nil {
					agg.Add(fmt.Errorf("marshal info_blocks_snapshot: %w", e))
					return
				}

				v := datatypes.JSON(raw)
				patch := map[string]any{"info_blocks_snapshot": v}
				if _, e = s.categoryRepo.UpdateFieldsByUUIDs([]uuid.UUID{id}, patch); e != nil {
					agg.Add(fmt.Errorf("update info_blocks_snapshot: %w", e))
					return
				}

				c.InfoBlocksSnapshot = v
			}
		}(category, category.UUID))
	}

	if len(category.InfoBlocksSnapshot) > 2 {
		var blocks []dto.InfoBlock
		if err := json.Unmarshal(category.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][categoryAggregate][InfoBlocks] id=%v: %v", category.ID, err)
		} else {
			interfaceBlocks := make([]contract.InfoBlock, 0, len(blocks))
			for _, block := range blocks {
				interfaceBlocks = append(interfaceBlocks, block)
			}

			category.InfoBlocks = interfaceBlocks
		}
	}
}

func (s *CategoryAggregateService) addGalleries(category *models.PostCategory, wg *sync.WaitGroup, agg *errutil.ErrUtil) {
	if category.GalleriesSnapshot == nil {
		service.SafeGo(wg, func(c *models.PostCategory, id uuid.UUID) func() {
			return func() {
				galleries := s.api.Gallery.GetForResourceUUID(id.String())
				if len(galleries) == 0 {
					return
				}

				raw, e := json.Marshal(dto.MapGalleries(galleries))
				if e != nil {
					agg.Add(fmt.Errorf("marshal galleries_snapshot: %w", e))
					return
				}

				v := datatypes.JSON(raw)
				patch := map[string]any{"galleries_snapshot": v}
				if _, e = s.categoryRepo.UpdateFieldsByUUIDs([]uuid.UUID{id}, patch); e != nil {
					agg.Add(fmt.Errorf("update galleries_snapshot: %w", e))
					return
				}

				c.GalleriesSnapshot = v
			}
		}(category, category.UUID))
	}

	if len(category.GalleriesSnapshot) > 2 {
		var galleries []dto.Gallery
		if err := json.Unmarshal(category.GalleriesSnapshot, &galleries); err != nil {
			logger.Errorf("[blog][categoryAggregate][Galleries] id=%v: %v", category.ID, err)
		} else {
			interfaceGalleries := make([]contract.Gallery, 0, len(galleries))
			for _, gallery := range galleries {
				interfaceGalleries = append(interfaceGalleries, gallery)
			}

			category.Galleries = interfaceGalleries
		}
	}
}
