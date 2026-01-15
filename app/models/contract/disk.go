package contract

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

// DiskService представляет сервис для чтения файлов с диска или из встроенных ресурсов
type DiskService interface {
	// ReadFile читает файл по указанному пути
	// Путь должен быть относительным от корня src (например: "services/i18n/locales/en.json")
	// Возвращает содержимое файла или ошибку, если файл не найден
	ReadFile(path string) ([]byte, error)

	ReadDir(pathString string) ([]fs.DirEntry, error)

	// ReadFileString читает файл и возвращает его содержимое как строку
	ReadFileString(path string) (string, error)

	// Exists проверяет существование файла
	Exists(path string) bool

	// GetTemplatesFS возвращает файловую систему для шаблонов
	GetTemplatesFS() fs.FS

	// GetStaticFS возвращает файловую систему для публичных файлов
	GetStaticFS() fs.FS

	// SetupStaticFiles настраивает роуты для статических файлов (/uploads и /static)
	SetupStaticFiles(router *gin.Engine) error
}
