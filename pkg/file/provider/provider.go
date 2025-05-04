package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/file/service"
	"mime/multipart"
)

type FileProvider interface {
	UploadFile(file *multipart.FileHeader, dist string) (path string, err error)
	UploadFiles(files []*multipart.FileHeader, dist string) (paths []string)
	DeleteFile(file string) error
	Received(files []string) error
	Exist(file string) bool
	RevisionReceived()
}

func NewFileProvider(
	uploadService *service.UploadService,
	fileService *service.Service,
	collectionService *service.CollectionService,
) FileProvider {
	return &provider{
		uploadService:     uploadService,
		fileService:       fileService,
		collectionService: collectionService,
	}
}

type provider struct {
	uploadService     *service.UploadService
	fileService       *service.Service
	collectionService *service.CollectionService
}

func (p *provider) UploadFile(file *multipart.FileHeader, dist string) (path string, err error) {
	return
}

func (p *provider) UploadFiles(files []*multipart.FileHeader, dist string) (paths []string) {
	return
}

func (p *provider) DeleteFile(file string) error {
	return p.fileService.Delete(file)
}

func (p *provider) Exist(file string) bool {
	return p.uploadService.Exist(file)
}

func (p *provider) Received(files []string) error { // @todo вернуть ошибки по файлам и удалить сущности при наличии ошибок
	return p.collectionService.Received(files)
}

func (p *provider) RevisionReceived() {
	err := p.collectionService.RevisionReceived()
	if err != nil {
		logger.Errorf("[FileProvider][RevisionReceived] Error: %v", err)
	}
}
