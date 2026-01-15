package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	http "github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/queue/job"
)

func (s *PostService) SaveFromRequest(form *http.PostRequest, found *models.Post, user contract.User) (model *models.Post, err error) {
	postForm := app.LoadStruct(&models.Post{}, form).(*models.Post)

	if found != nil {
		model, err = s.Update(postForm, found)
	} else {
		model, err = s.Create(postForm, user)
	}

	if err != nil {
		return model, err
	}

	if len(form.Galleries) > 0 {
		interfaceSlice := make([]any, len(form.Galleries))
		for i, gallery := range form.Galleries {
			interfaceSlice[i] = gallery
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

	if len(form.Tags) > 0 {
		tags, err := s.tagCollectionService.SyncTags(model.UUID, form.Tags)
		if err != nil {
			return nil, err
		}

		model.PostTags = tags
	}

	return model, nil
}

func (s *PostService) Create(model *models.Post, user contract.User) (*models.Post, error) {
	id := user.GetID()
	model.UserID = &id
	model.Alias = s.generateAlias(model)

	if err := s.postRepo.Create(model); err != nil {
		return nil, err
	}

	if err := s.receivedImage(model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *PostService) Update(model, found *models.Post) (*models.Post, error) {
	model.ID = found.ID
	model.UUID = found.UUID
	model.UserID = found.UserID

	if model.Alias != found.Alias {
		model.Alias = s.generateAlias(model)
	}

	if err := s.postRepo.Update(model); err != nil {
		return nil, err
	}

	if model.Image != nil && (found.Image == nil || *model.Image != *found.Image) {
		if err := s.receivedImage(model); err != nil {
			return nil, err
		}
	}

	s.queue.Enqueue(job.NewPostJob(model, "update"), 0)

	return model, nil
}
