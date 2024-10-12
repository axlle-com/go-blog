package service

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	common "github.com/axlle-com/blog/pkg/common/service"
	gallery "github.com/axlle-com/blog/pkg/gallery/service"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/models"
	user "github.com/axlle-com/blog/pkg/user/models"
)

func PostCreate(form *http.PostRequest, u *user.User) (*models.Post, error) {
	post := common.LoadStruct(&models.Post{}, form).(*models.Post)
	post.UserID = &u.ID
	if err := models.PostRepo().Create(post); err != nil {
		return nil, err
	}

	if len(form.Galleries) > 0 {
		slice := make([]contracts.Gallery, len(form.Galleries), len(form.Galleries))
		for _, i := range form.Galleries {
			g, _ := gallery.SaveGallery(i)
			err := g.Attach(post)
			if err != nil {
				continue
			}
			slice = append(slice, g)
		}
		post.Galleries = slice
	}
	return post, nil
}
