package service

import (
	"github.com/axlle-com/blog/pkg/file/models"
	"github.com/axlle-com/blog/pkg/file/repository"
)

type Service struct {
	fileRepo repository.FileRepository
}

func NewService(
	fileRepo repository.FileRepository,
) *Service {
	return &Service{
		fileRepo: fileRepo,
	}
}

func (s *Service) Create(file *models.File) error {
	return s.fileRepo.Create(file)
}

func (s *Service) Delete(file string) error {
	model, err := s.fileRepo.GetByFile(file)
	if err != nil {
		return err
	}

	return s.fileRepo.Delete(model.ID)
}

func (s *Service) Destroy(file string) error {
	model, err := s.fileRepo.GetByFile(file)
	if err != nil {
		return err
	}

	return s.fileRepo.Destroy(model.ID)
}

func (s *Service) Received(file string) error {
	return s.fileRepo.Received([]string{file})
}
