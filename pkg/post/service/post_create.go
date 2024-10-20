package service

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	common "github.com/axlle-com/blog/pkg/common/service"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/models"
	user "github.com/axlle-com/blog/pkg/user/models"
)

func PostSave(form *http.PostRequest, u *user.User) (*models.Post, error) {
	post := common.LoadStruct(&models.Post{}, form).(*models.Post)
	repo := models.PostRepo()
	post.UserID = &u.ID

	if post.ID == 0 {
		if err := repo.Create(post); err != nil {
			return nil, err
		}
	} else {
		if err := repo.Update(post); err != nil {
			return nil, err
		}
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts.Gallery, len(form.Galleries), len(form.Galleries))
		for idx, gRequest := range form.Galleries {
			if gRequest == nil {
				continue
			}
			g, err := gallery.Provider().SaveFromForm(gRequest)
			if err != nil {
				continue
			}
			err = g.Attach(post)
			if err != nil {
				continue
			}
			slice[idx] = g
		}
		post.Galleries = slice
	}
	return post, nil
}
