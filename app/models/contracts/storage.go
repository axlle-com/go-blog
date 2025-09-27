package contracts

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	// Save сохраняет поток в целевой folder с именем <uuid><ext> и возвращает ПОЛНЫЙ http(s) URL
	Save(fileHeader *multipart.FileHeader, folder, fileName string) (url string, err error)
	Destroy(urlOrPath string) error
	Exists(urlOrPath string) bool
	// Optional: прямое сохранение из io.Reader (на будущее)
	SaveReader(r io.Reader, size int64, folder, filename string, contentType string) (url string, err error)
}
