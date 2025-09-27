package service

import (
	"github.com/axlle-com/blog/pkg/file/models"
	"github.com/axlle-com/blog/pkg/file/repository"
)

type FileService struct {
	fileRepo repository.FileRepository
}

func NewFileService(
	fileRepo repository.FileRepository,
) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}

func (s *FileService) Create(file *models.File) error {
	return s.fileRepo.Create(file)
}

func (s *FileService) Delete(file string) error {
	model, err := s.fileRepo.GetByFile(file)
	if err != nil {
		return err
	}
	return s.fileRepo.Delete(model.ID)
}

func (s *FileService) Destroy(file string) error {
	model, err := s.fileRepo.GetByFile(file)
	if err != nil {
		return err
	}
	return s.fileRepo.Destroy(model.ID)
}

func (s *FileService) Received(file string) error {
	return s.fileRepo.Received([]string{file})
}
