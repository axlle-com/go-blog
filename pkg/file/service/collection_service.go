package service

import (
	"fmt"
	"os"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/pkg/file/repository"
)

type CollectionService struct {
	fileRepo      repository.FileRepository
	service       *Service
	uploadService *UploadService
}

func NewCollectionService(
	fileRepo repository.FileRepository,
	fileService *Service,
	uploadService *UploadService,
) *CollectionService {
	return &CollectionService{
		fileRepo:      fileRepo,
		service:       fileService,
		uploadService: uploadService,
	}
}

func (s *CollectionService) Received(files []string) error {
	return s.fileRepo.Received(files)
}

func (s *CollectionService) RevisionReceived() error {
	params := map[string]any{
		"received_at": nil,
	}

	byParams, err := s.fileRepo.GetByParams(params, true)
	if err != nil {
		return err
	}

	newErr := errutil.New()

	for _, file := range byParams {
		if err := s.uploadService.DestroyFile(file.File); err != nil {
			if !os.IsNotExist(err) {
				newErr.Add(fmt.Errorf("DestroyFile(%s): %v", file.File, err))
			}
			continue
		}

		if err := s.fileRepo.Destroy(file.ID); err != nil {
			newErr.Add(fmt.Errorf("DestroyRecord(%d): %v", file.ID, err))
		}
	}

	return newErr.Error()
}
