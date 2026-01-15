package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/app/service"
	app "github.com/axlle-com/blog/app/service/struct"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo          repository.CategoryRepository
	postCollectionService *PostCollectionService
	api                   *api.Api
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	api *api.Api,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		api:          api,
	}
}

func (s *CategoryService) SetPostCollectionService(postCollectionService *PostCollectionService) {
	s.postCollectionService = postCollectionService
}

func (s *CategoryService) GetAggregateByID(id uint) (*models.PostCategory, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.Aggregate(category)
}

func (s *CategoryService) Aggregate(category *models.PostCategory) (*models.PostCategory, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)

	wg.Add(2)

	go func() {
		defer wg.Done()
		galleries = s.api.Gallery.GetForResourceUUID(category.UUID.String())
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(category.UUID.String())
	}()

	wg.Wait()

	category.Galleries = galleries
	category.InfoBlocks = infoBlocks

	return category, nil
}

func (s *CategoryService) SaveFromRequest(
	form *http.CategoryRequest,
	found *models.PostCategory,
	user contract.User,
) (model *models.PostCategory, err error) {
	categoryForm := app.LoadStruct(&models.PostCategory{}, form).(*models.PostCategory)

	if found == nil {
		model, err = s.Create(categoryForm, user)
	} else {
		model, err = s.Update(categoryForm, found, user)
	}

	if err != nil {
		return
	}

	if len(form.Galleries) > 0 {
		interfaceSlice := make([]any, len(form.Galleries))
		for i, gall := range form.Galleries {
			interfaceSlice[i] = gall
		}

		slice, err := s.api.Gallery.SaveFormBatch(interfaceSlice, model)
		if err != nil {
			logger.Error(err)
		}
		model.Galleries = slice
	}

	if len(form.InfoBlocks) > 0 {
		interfaceSlice := make([]any, len(form.InfoBlocks))
		for i, block := range form.InfoBlocks {
			interfaceSlice[i] = block
		}

		slice, err := s.api.InfoBlock.CreateRelationFormBatch(interfaceSlice, model.UUID.String())
		if err != nil {
			logger.Error(err)
		}
		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *CategoryService) GetByID(id uint) (*models.PostCategory, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) Delete(category *models.PostCategory) error {
	err := s.api.Gallery.DetachResource(category)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(category)
}

func (s *CategoryService) Create(model *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	id := user.GetID()
	model.UserID = &id
	model.Alias = s.generateAlias(model)

	if err := s.categoryRepo.Create(model); err != nil {
		return nil, err
	}

	if err := s.receivedImage(model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *CategoryService) Update(model, found *models.PostCategory, user contract.User) (*models.PostCategory, error) {
	model.ID = found.ID
	model.UUID = found.UUID
	model.UserID = found.UserID

	if model.Alias != found.Alias {
		model.Alias = s.generateAlias(model)
	}

	tx := s.categoryRepo.Tx()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.categoryRepo.WithTx(tx).Update(model, found); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if model.Image != nil && (found.Image == nil || *model.Image != *found.Image) {
		if err := s.receivedImage(model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (s *CategoryService) generateAlias(category *models.PostCategory) string {
	var aliasStr string
	if category.Alias == "" {
		aliasStr = category.Title
	} else {
		aliasStr = category.Alias
	}

	return s.api.Alias.Generate(category, aliasStr)
}

func (s *CategoryService) DeleteImageFile(category *models.PostCategory) error {
	if category.Image == nil {
		return nil
	}

	err := s.api.File.DeleteFile(*category.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	category.Image = nil

	return nil
}

func (s *CategoryService) View(
	category *models.PostCategory,
	paginator contract.Paginator,
	filter *models.PostFilter,
) (*models.PostCategory, error) {
	var wg sync.WaitGroup
	agg := errutil.New()

	// --- Posts with pagination ---
	if s.postCollectionService != nil && paginator != nil {
		service.SafeGo(&wg, func(c *models.PostCategory, p contract.Paginator, f *models.PostFilter) func() {
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

	// --- InfoBlocks snapshot ---
	if category.InfoBlocksSnapshot == nil {
		service.SafeGo(&wg, func(c *models.PostCategory, id uuid.UUID) func() {
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

	// --- Galleries snapshot ---
	if category.GalleriesSnapshot == nil {
		service.SafeGo(&wg, func(c *models.PostCategory, id uuid.UUID) func() {
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

	// --- template ---
	if category.TemplateID != nil {
		service.SafeGo(&wg, func(c *models.PostCategory) func() {
			return func() {
				tpl, e := s.api.Template.GetByID(*c.TemplateID)
				if e != nil {
					agg.Add(fmt.Errorf("get template: %w", e))
					return
				}
				c.Template = tpl
			}
		}(category))
	}

	wg.Wait()

	return category, agg.Error()
}

// @todo переделать на фильтр везде
func (s *CategoryService) FindByParam(field string, value any) (*models.PostCategory, error) {
	return s.categoryRepo.FindByParam(field, value)
}

func (s *CategoryService) receivedImage(category *models.PostCategory) error {
	if category.Image != nil && *category.Image != "" {
		return s.api.File.Received([]string{*category.Image})
	}

	return nil
}
