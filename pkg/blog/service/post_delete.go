package service

import (
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/queue/job"
)

func (s *PostService) PostDelete(post *models.Post) error {
	err := s.api.Gallery.DetachResource(post)
	if err != nil {
		return err
	}

	err = s.api.InfoBlock.DetachResourceUUID(post.UUID.String())
	if err != nil {
		return err
	}
	if err := s.postRepo.Delete(post); err != nil {
		return err
	}

	s.queue.Enqueue(job.NewPostJob(post, "delete"), 0)

	return nil
}
