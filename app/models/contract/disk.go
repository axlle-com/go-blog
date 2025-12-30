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

	// ReadFileString читает файл и возвращает его содержимое как строку
	ReadFileString(path string) (string, error)

	// Exists проверяет существование файла
	Exists(path string) bool

	// GetTemplatesFS возвращает файловую систему для шаблонов
	GetTemplatesFS() fs.FS

	// GetPublicFS возвращает файловую систему для публичных файлов
	GetPublicFS() fs.FS

	// SetupStaticFiles настраивает роуты для статических файлов (/uploads и /public)
	SetupStaticFiles(router *gin.Engine) error
}
