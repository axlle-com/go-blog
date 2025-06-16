package service

import (
	"errors"
	"github.com/axlle-com/blog/app/logger"
	"gorm.io/gorm"

	"github.com/axlle-com/blog/pkg/blog/models"
)

func (s *PostService) DeleteImageFile(post *models.Post) error {
	if post.Image == nil {
		return nil
	}
	err := s.fileProvider.DeleteFile(*post.Image)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		logger.Errorf("[DeleteImageFile] Error: %v", err)
	}
	post.Image = nil
	return nil
}
