package service

import (
	"errors"
	"gorm.io/gorm"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/post/models"
)

func (s *PostService) DeleteImageFile(post *models.Post) error {
	if post.Image == nil {
		return nil
	}
	err := s.fileProvider.DeleteFile(*post.Image)
	if err != nil {
		logger.Errorf("[DeleteImageFile] Error: %v", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	post.Image = nil
	return nil
}
