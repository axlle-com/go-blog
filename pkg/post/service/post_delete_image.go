package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/post/models"
)

func (s *Service) DeleteImageFile(post *models.Post) error {
	if post.Image == nil {
		return errors.New("image is nil")
	}
	err := s.fileProvider.DeleteFile(*post.Image)
	if err != nil {
		return err
	}
	post.Image = nil
	return nil
}
