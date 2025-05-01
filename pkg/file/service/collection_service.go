package service

import (
	"github.com/axlle-com/blog/pkg/file/repository"
)

type CollectionService struct {
	fileRepo repository.FileRepository
}

func NewCollectionService(
	fileRepo repository.FileRepository,
) *CollectionService {
	return &CollectionService{
		fileRepo: fileRepo,
	}
}

func (s *CollectionService) Received(files []string) error {
	return s.fileRepo.Received(files)
}
