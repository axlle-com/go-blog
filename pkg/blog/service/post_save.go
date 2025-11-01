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

		slice, err := s.galleryProvider.SaveFormBatch(interfaceSlice, model)
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

		slice, err := s.infoBlockProvider.SaveFormBatch(interfaceSlice, model.UUID.String())
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

func (s *PostService) Create(post *models.Post, user contract.User) (*models.Post, error) {
	id := user.GetID()
	post.UserID = &id
	post.Alias = s.generateAlias(post)

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	if err := s.receivedImage(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Update(post, found *models.Post) (*models.Post, error) {
	post.ID = found.ID
	post.UUID = found.UUID

	if post.Alias != found.Alias {
		post.Alias = s.generateAlias(post)
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	if post.Image != nil && found.Image != nil && *post.Image != *found.Image {
		if err := s.receivedImage(post); err != nil {
			return nil, err
		}
	}

	s.queue.Enqueue(job.NewPostJob(post, "update"), 0)

	return post, nil
}
