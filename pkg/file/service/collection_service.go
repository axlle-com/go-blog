package service

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/file/repository"
	"os"
	"strings"
)

type CollectionService struct {
	fileRepo      repository.FileRepository
	fileService   *Service
	uploadService *UploadService
}

func NewCollectionService(
	fileRepo repository.FileRepository,
	fileService *Service,
	uploadService *UploadService,
) *CollectionService {
	return &CollectionService{
		fileRepo:      fileRepo,
		fileService:   fileService,
		uploadService: uploadService,
	}
}

func (s *CollectionService) Received(files []string) error {
	return s.fileRepo.Received(files)
}

func (s *CollectionService) RevisionReceived() error {
	params := map[string]any{
		"is_received": false,
	}

	byParams, err := s.fileRepo.GetByParams(params, true)
	if err != nil {
		return err
	}

	var errs []string

	for _, file := range byParams {
		if err := s.uploadService.DestroyFile(file.File); err != nil {
			if !os.IsNotExist(err) {
				errs = append(errs, fmt.Sprintf("DestroyFile(%s): %v", file.File, err))
			}
			continue
		}

		if err := s.fileRepo.Destroy(file.ID); err != nil {
			errs = append(errs, fmt.Sprintf("Destroy(record %d): %v", file.ID, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
