package service

import (
	"errors"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	"gorm.io/gorm"
)

func (s *PostService) DeleteImageFile(post *models.Post) error {
	if post.Image == nil {
		return nil
	}
	err := s.api.File.DeleteFile(*post.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		logger.Errorf("[DeleteImageFile] Error: %v", err)
	}
	post.Image = nil
	return nil
}
