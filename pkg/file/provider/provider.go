package provider

import (
	"github.com/axlle-com/blog/pkg/file/service"
	"mime/multipart"
)

type FileProvider interface {
	UploadFile(file *multipart.FileHeader, dist string) (path string, err error)
	UploadFiles(files []*multipart.FileHeader, dist string) (paths []string)
	DeleteFile(file string) error
	Received(files []string) error
}

func NewProvider(
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

func (p *provider) Received(files []string) error {
	return p.collectionService.Received(files)
}
