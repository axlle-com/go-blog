package disk

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/src"
	"github.com/gin-gonic/gin"
)

type diskService struct {
	config      contract.Config
	servicesFS  embed.FS
	publicFS    embed.FS
	templatesFS embed.FS

	servicesRoot  string // "services"
	publicRoot    string // "public"
	templatesRoot string // "templates"
}

func NewDiskService(cfg contract.Config) contract.DiskService {
	return &diskService{
		config:        cfg,
		servicesFS:    src.ServicesFS,
		publicFS:      src.PublicFS,
		templatesFS:   src.TemplatesFS,
		servicesRoot:  "services",
		publicRoot:    "public",
		templatesRoot: "templates",
	}
}

// normalizePath приводит путь к виду, подходящему для fs.FS (всегда "/"),
// убирает ведущие "/" и чистит "."/"..".
func normalizePath(pathString string) string {
	pathString = strings.TrimSpace(pathString)
	pathString = filepath.ToSlash(pathString)
	pathString = strings.TrimLeft(pathString, "/")
	pathString = path.Clean(pathString)
	if pathString == "." {
		return ""
	}
	return pathString
}

// resolvePath ожидает, что p уже normalizePath(...).
// Возвращает FS и полный путь внутри неё.
func (d *diskService) resolvePath(pathString string) (fs.FS, string) {
	switch {
	case pathString == "":
		return nil, ""
	case strings.HasPrefix(pathString, d.templatesRoot+"/") || pathString == d.templatesRoot:
		return d.templatesFS, pathString
	case strings.HasPrefix(pathString, d.publicRoot+"/") || pathString == d.publicRoot:
		return d.publicFS, pathString
	case strings.HasPrefix(pathString, d.servicesRoot+"/") || pathString == d.servicesRoot:
		return d.servicesFS, pathString
	default:
		// по умолчанию считаем, что это services/<pathString>
		return d.servicesFS, path.Join(d.servicesRoot, pathString)
	}
}

func (d *diskService) ReadFile(pathString string) ([]byte, error) {
	pathString = normalizePath(pathString)
	if pathString == "" {
		return nil, fs.ErrNotExist
	}

	fsys, fullPath := d.resolvePath(pathString)
	if fsys != nil {
		data, err := fs.ReadFile(fsys, fullPath)
		if err == nil {
			return data, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	// Если не найдено в embed, пытаемся прочитать с диска
	diskPath := d.config.SrcFolderBuilder(pathString)

	// Проверяем существование файла на диске
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		return nil, err
	}

	// Читаем с диска
	data, err := os.ReadFile(diskPath)
	if err != nil {
		logger.Errorf("[disk][ReadFile] error reading file from disk: %s, error: %v", diskPath, err)
		return nil, err
	}

	logger.Debugf("[disk][ReadFile] read file from disk: %s", diskPath)

	return data, nil
}

func (d *diskService) ReadFileString(p string) (string, error) {
	b, err := d.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (d *diskService) Exists(p string) bool {
	p = normalizePath(p)
	if p == "" {
		return false
	}

	fsys, fullPath := d.resolvePath(p)
	if fsys != nil {
		if _, err := fs.Stat(fsys, fullPath); err == nil {
			return true
		} else if !errors.Is(err, fs.ErrNotExist) {
			// неожиданные ошибки игнорируем (считаем, что "нет")
			return false
		}
	}

	// disk fallback только в local (dev)
	if d.config != nil && d.config.IsLocal() {
		diskPath := d.config.SrcFolderBuilder(p)
		if _, err := os.Stat(diskPath); err == nil {
			return true
		}
	}

	return false
}

func (d *diskService) GetTemplatesFS() fs.FS {
	return d.templatesFS
}

func (d *diskService) GetPublicFS() fs.FS {
	return d.publicFS
}

func (d *diskService) SetupStaticFiles(router *gin.Engine) error {
	if router == nil {
		return nil
	}

	// Добавляем заголовки Cache-Control для статических файлов
	router.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/public/") || strings.HasPrefix(p, "/uploads/") || p == "/favicon.ico" {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	})

	// /uploads/* - всегда с диска (пользовательские загруженные файлы)
	// Физически файлы в data/uploads, но URL остаются /uploads/
	router.Static("/uploads", d.config.DataFolder(d.config.UploadPath()))

	// /public/* - только из embed (статические файлы)
	publicSub, err := fs.Sub(d.publicFS, "public")
	if err != nil {
		logger.Fatalf("[disk][SetupStaticFiles] fs.Sub(public) failed: %v", err)
		return err
	}
	router.StaticFS("/public", http.FS(publicSub))

	// favicon.ico из embed
	router.GET("/favicon.ico", func(c *gin.Context) {
		b, err := fs.ReadFile(d.publicFS, "public/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "image/x-icon", b)
	})

	return nil
}
