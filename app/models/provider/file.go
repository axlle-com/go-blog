package provider

import "mime/multipart"

type FileProvider interface {
	UploadFile(file *multipart.FileHeader, dist string) (path string, err error)
	UploadFiles(files []*multipart.FileHeader, dist string) (paths []string)
	DeleteFile(file string) error
	Received(files []string) error
	Exist(file string) bool
	RevisionReceived()
}
