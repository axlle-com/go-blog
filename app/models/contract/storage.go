package contract

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	Save(fileHeader *multipart.FileHeader, folder, fileName string) (url string, err error)
	Destroy(urlOrPath string) error
	Exists(urlOrPath string) bool
	SaveReader(r io.Reader, size int64, folder, filename string, contentType string) (url string, err error)
}
