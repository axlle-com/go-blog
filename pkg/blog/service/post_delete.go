package service

import (
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/queue"
)

func (s *PostService) PostDelete(post *models.Post) error {
	err := s.galleryProvider.DetachResource(post)
	if err != nil {
		return err
	}

	err = s.infoBlockProvider.DetachResource(post)
	if err != nil {
		return err
	}
	if err := s.postRepo.Delete(post); err != nil {
		return err
	}

	s.queue.Enqueue(queue.NewPostJob(post, "delete"), 0)

	return nil
}
