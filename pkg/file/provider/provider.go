package provider

import (
	"github.com/axlle-com/blog/pkg/file"
	"mime/multipart"
)

type FileProvider interface {
	UploadFile(file *multipart.FileHeader, dist string) (path string, err error)
	UploadFiles(files []*multipart.FileHeader, dist string) (paths []string)
	DeleteFile(file string) error
}

func NewProvider(
	service *file.Service,
) FileProvider {
	return &provider{
		service: service,
	}
}

type provider struct {
	service *file.Service
}

func (p *provider) UploadFile(file *multipart.FileHeader, dist string) (path string, err error) {
	return p.service.SaveUploadedFile(file, dist)
}

func (p *provider) UploadFiles(files []*multipart.FileHeader, dist string) (paths []string) {
	return p.service.SaveUploadedFiles(files, dist)
}

func (p *provider) DeleteFile(file string) error {
	return p.service.DeleteFile(file)
}
