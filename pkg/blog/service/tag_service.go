package service

import (
	"errors"
	"sync"
	"unicode/utf8"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service"
	app "github.com/axlle-com/blog/app/service/struct"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"gorm.io/gorm"
)

type TagService struct {
	tagRepo      repository.PostTagRepository
	resourceRepo repository.PostTagResourceRepository
	api          *api.Api
}

func NewTagService(
	postTagRepo repository.PostTagRepository,
	resourceRepo repository.PostTagResourceRepository,
	api *api.Api,
) *TagService {
	return &TagService{
		tagRepo:      postTagRepo,
		resourceRepo: resourceRepo,
		api:          api,
	}
}

func (s *TagService) GetByID(id uint) (*models.PostTag, error) {
	return s.tagRepo.GetByID(id)
}

func (s *TagService) Aggregate(model *models.PostTag) (*models.PostTag, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)

	service.SafeGo(&wg, func() {
		galleries = s.api.Gallery.GetForResourceUUID(model.UUID.String())
	})

	service.SafeGo(&wg, func() {
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(model.UUID.String())
	})

	service.SafeGo(&wg, func() {
		if model.TemplateName == "" {
			return
		}

		tpl, err := s.api.Template.GetByName(model.TemplateName)
		if err != nil {
			logger.Errorf("[TagService] Error: %+v", err)
			return
		}
		model.Template = tpl
	})

	wg.Wait()

	model.Galleries = galleries
	model.InfoBlocks = infoBlocks

	return model, nil
}

func (s *TagService) SaveFromRequest(form *http.TagRequest, found *models.PostTag, user contract.User) (model *models.PostTag, err error) {
	tagForm := app.LoadStruct(&models.PostTag{}, form).(*models.PostTag)

	if found != nil {
		model, err = s.Update(tagForm, found)
	} else {
		model, err = s.Create(tagForm)
	}

	if err != nil {
		return model, err
	}

	if len(form.Galleries) > 0 {
		interfaceSlice := make([]any, len(form.Galleries))
		for i, galleryRequest := range form.Galleries {
			interfaceSlice[i] = galleryRequest
		}

		slice, err := s.api.Gallery.SaveFormBatch(interfaceSlice, model)
		if err != nil {
			logger.Errorf("[blog][TagService][SaveFromRequest] Error: %+v", err)
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
			logger.Errorf("[blog][TagService][SaveFromRequest] Error: %+v", err)
		}
		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *TagService) Create(model *models.PostTag) (*models.PostTag, error) {
	model.Name = s.trimName(model.Name, 10)
	model.Alias = s.generateAlias(model)

	if err := s.tagRepo.Create(model); err != nil {
		return nil, err
	}

	if err := s.receivedImage(model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *TagService) Update(model, found *models.PostTag) (*models.PostTag, error) {
	model.Name = s.trimName(model.Name, 10)

	if model.Alias != found.Alias {
		model.Alias = s.generateAlias(model)
	}

	if err := s.tagRepo.Update(model); err != nil {
		return nil, err
	}

	if model.Image != nil && (found.Image == nil || *model.Image != *found.Image) {
		if err := s.receivedImage(model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (s *TagService) Attach(resource contract.Resource, postTag contract.PostTag) error {
	hasRepo, err := s.resourceRepo.GetByParams(resource.GetUUID(), postTag.GetID())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.PostTagHasResource{
				ResourceUUID: resource.GetUUID(),
				PostTagID:    postTag.GetID(),
			},
		)
		return err
	}
	return nil
}

func (s *TagService) DeleteTags(postTags []*models.PostTag) (err error) {
	var ids []uint
	for _, postTag := range postTags {
		ids = append(ids, postTag.ID)
	}

	if len(ids) > 0 {
		if err = s.tagRepo.DeleteByIDs(ids); err != nil {
			return err
		}

		if err = s.resourceRepo.DeleteByIDs(ids); err != nil {
			return err
		}
	}

	return nil
}

func (s *TagService) DeleteImageFile(tag *models.PostTag) error {
	if tag.Image == nil {
		return nil
	}

	err := s.api.File.DeleteFile(*tag.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	tag.Image = nil

	return nil
}

func (s *TagService) receivedImage(model *models.PostTag) error {
	if model.Image != nil && *model.Image != "" {
		return s.api.File.Received([]string{*model.Image})
	}

	return nil
}

func (s *TagService) generateAlias(model *models.PostTag) string {
	var newAlias string

	if model.Alias == "" {
		if model.Name == "" && model.Title != nil {
			newAlias = *model.Title
		} else {
			newAlias = model.Name
		}
	} else {
		newAlias = model.Alias
	}

	return s.api.Alias.Generate(model, newAlias)
}

func (s *TagService) trimName(name string, max int) string {
	if max <= 0 || name == "" {
		return ""
	}

	if utf8.RuneCountInString(name) <= max {
		return name
	}

	r := []rune(name)

	return string(r[:max])
}
