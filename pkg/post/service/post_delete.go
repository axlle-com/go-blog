package service

import (
	"github.com/axlle-com/blog/pkg/post/models"
)

func (s *Service) PostDelete(post *models.Post) error {
	err := s.galleryProvider.DeleteForResource(post)
	if err != nil {
		return err
	}
	if err := s.postRepo.Delete(post); err != nil {
		return err
	}
	return nil
}
