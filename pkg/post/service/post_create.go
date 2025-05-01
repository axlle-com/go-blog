package service

import (
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/models"
)

func (s *PostService) SaveFromRequest(form *http.PostRequest, user contracts2.User) (*models.Post, error) {
	postForm := app.LoadStruct(&models.Post{}, form).(*models.Post)
	model, err := s.Save(postForm, user)
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

		slice, err := s.infoBlockProvider.SaveFormBatch(interfaceSlice, model)
		if err != nil {
			logger.Error(err)
		}
		model.InfoBlocks = slice
	}

	return model, nil
}

func (s *PostService) Save(post *models.Post, user contracts2.User) (*models.Post, error) {
	id := user.GetID()
	post.UserID = &id
	var alias string
	if post.Alias == "" {
		alias = post.Title
	} else {
		alias = post.Alias
	}

	post.Alias = s.aliasProvider.Generate(post, alias)
	if post.ID == 0 {
		if err := s.postRepo.Create(post); err != nil {
			return nil, err
		}
	} else {
		if err := s.postRepo.Update(post); err != nil {
			return nil, err
		}
	}

	if post.Image != nil && *post.Image != "" {
		err := s.fileProvider.Received([]string{*post.Image})
		if err != nil {
			return post, err
		}
	}

	return post, nil
}
