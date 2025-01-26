package service

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	app "github.com/axlle-com/blog/pkg/app/service"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/models"
	user "github.com/axlle-com/blog/pkg/user/models"
)

func (s *Service) SaveFromRequest(form *http.PostRequest, user *user.User) (*models.Post, error) {
	postForm := app.LoadStruct(&models.Post{}, form).(*models.Post)
	post, err := s.Save(postForm, user)
	if err != nil {
		return nil, err
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts.Gallery, 0)
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
	return post, nil
}

func (s *Service) Save(post *models.Post, user *user.User) (*models.Post, error) {
	post.UserID = &user.ID
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
