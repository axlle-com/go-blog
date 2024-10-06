package service

import (
	common "github.com/axlle-com/blog/pkg/common/service"
	http "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/axlle-com/blog/pkg/post/models"
	user "github.com/axlle-com/blog/pkg/user/models"
)

func PostCreate(form *http.PostRequest, u *user.User) (*models.Post, error) {
	post := common.LoadFromRequest(&models.Post{}, form).(*models.Post)
	post.UserID = &u.ID
	if err := models.PostRepo().Create(post); err != nil {
		return nil, err
	}

	return post, nil
}
