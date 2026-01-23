package contract

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

// DiskService представляет сервис для чтения файлов с диска или из встроенных ресурсов
type DiskService interface {
	ReadFile(path string) ([]byte, error)
	ReadDir(pathString string) ([]fs.DirEntry, error)
	ReadFileString(path string) (string, error)
	Exists(path string) bool
	GetTemplatesFS() fs.FS
	GetStaticFS() fs.FS
	SetupStaticFiles(router *gin.Engine) error
}
