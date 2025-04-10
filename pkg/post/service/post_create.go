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
	post, err := s.Save(postForm, user)
	if err != nil {
		return nil, err
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts2.Gallery, 0)
		for _, gRequest := range form.Galleries {
			if gRequest == nil {
				continue
			}

			g, err := s.galleryProvider.SaveFromForm(gRequest, post)
			if err != nil || g == nil {
				continue
			}
			slice = append(slice, g)
		}
		post.Galleries = slice
	}

	if len(form.InfoBlocks) > 0 {
		interfaceSlice := make([]any, len(form.InfoBlocks))
		for i, block := range form.InfoBlocks {
			interfaceSlice[i] = block
		}

		slice, err := s.infoBlockProvider.SaveFormBatch(interfaceSlice, post)
		if err != nil {
			logger.Error(err)
		}
		post.InfoBlocks = slice
	}

	return post, nil
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

	return post, nil
}
