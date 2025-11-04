package service

import (
	"errors"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
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

func (s *TagService) Aggregate(post *models.PostTag) (*models.PostTag, error) {
	var wg sync.WaitGroup

	var galleries = make([]contract.Gallery, 0)
	var infoBlocks = make([]contract.InfoBlock, 0)

	wg.Add(2)

	go func() {
		defer wg.Done()
		galleries = s.api.Gallery.GetForResourceUUID(post.UUID.String())
	}()

	go func() {
		defer wg.Done()
		infoBlocks = s.api.InfoBlock.GetForResourceUUID(post.UUID.String())
	}()

	wg.Wait()

	post.Galleries = galleries
	post.InfoBlocks = infoBlocks

	return post, nil
}

func (s *TagService) Create(postTag *models.PostTag) (*models.PostTag, error) {
	postTag.Alias = s.api.Alias.Generate(postTag, postTag.Name)
	if err := s.tagRepo.Create(postTag); err != nil {
		return nil, err
	}

	return postTag, nil
}

func (s *TagService) Update(postTag *models.PostTag) (*models.PostTag, error) {
	if err := s.tagRepo.Update(postTag); err != nil {
		return nil, err
	}

	return postTag, nil
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

func (s *TagService) SaveFromRequest(form *http.TagRequest, user contract.User) (*models.PostTag, error) {
	tagForm := app.LoadStruct(&models.PostTag{}, form).(*models.PostTag)
	model, err := s.Save(tagForm, user)
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
			logger.Error(err)
		}
		model.Galleries = slice
	}

	if len(form.InfoBlocks) > 0 {
		interfaceSlice := make([]any, len(form.InfoBlocks))
		for i, block := range form.InfoBlocks {
			interfaceSlice[i] = block
		}

		slice, err := s.api.InfoBlock.SaveFormBatch(interfaceSlice, model.UUID.String())
		if err != nil {
			logger.Error(err)
		}
		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *TagService) Save(tag *models.PostTag, user contract.User) (*models.PostTag, error) {
	var newAlias string
	if tag.Alias == "" {
		newAlias = *tag.Title
	} else {
		newAlias = tag.Alias
	}

	tag.Alias = s.api.Alias.Generate(tag, newAlias)
	if tag.ID == 0 {
		if err := s.tagRepo.Create(tag); err != nil {
			return nil, err
		}
	} else {
		if err := s.tagRepo.Update(tag); err != nil {
			return nil, err
		}
	}

	if tag.Image != nil && *tag.Image != "" {
		err := s.api.File.Received([]string{*tag.Image})
		if err != nil {
			return tag, err
		}
	}

	return tag, nil
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
		logger.Errorf("[DeleteImageFile] Error: %v", err)
	}
	tag.Image = nil
	return nil
}
