package disk

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/src"
	"github.com/gin-gonic/gin"
)

type diskService struct {
	config      contract.Config
	servicesFS  embed.FS
	staticFS    embed.FS
	templatesFS embed.FS

	servicesRoot  string // "services"
	staticRoot    string // "static"
	templatesRoot string // "templates"
}

func NewDiskService(cfg contract.Config) contract.DiskService {
	return &diskService{
		config:        cfg,
		servicesFS:    src.ServicesFS,
		staticFS:      src.StaticFS,
		templatesFS:   src.TemplatesFS,
		servicesRoot:  "services",
		staticRoot:    "static",
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

func (d *diskService) resolvePath(pathString string) (fs.FS, string) {
	switch {
	case strings.HasPrefix(pathString, d.templatesRoot+"/") || pathString == d.templatesRoot:
		return d.templatesFS, pathString
	case strings.HasPrefix(pathString, d.staticRoot+"/") || pathString == d.staticRoot:
		return d.staticFS, pathString
	case strings.HasPrefix(pathString, d.servicesRoot+"/") || pathString == d.servicesRoot:
		return d.servicesFS, pathString
	default:
		return nil, ""
	}
}

func (d *diskService) ReadDir(pathString string) ([]fs.DirEntry, error) {
	pathString = normalizePath(pathString)
	if pathString == "" {
		return nil, fs.ErrNotExist
	}

	merged := map[string]fs.DirEntry{}

	// 1) Пробуем embed
	if fsys, fullPath := d.resolvePath(pathString); fsys != nil {
		if entries, err := fs.ReadDir(fsys, fullPath); err == nil {
			for _, e := range entries {
				merged[e.Name()] = e
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	// 2) Пробуем диск
	diskPath := d.config.DataFolder(pathString)
	if entries, err := os.ReadDir(diskPath); err == nil {
		for _, e := range entries {
			merged[e.Name()] = e
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if len(merged) == 0 {
		return nil, fs.ErrNotExist
	}

	names := make([]string, 0, len(merged))
	for name := range merged {
		names = append(names, name)
	}

	sort.Strings(names)

	out := make([]fs.DirEntry, 0, len(names))
	for _, name := range names {
		out = append(out, merged[name])
	}
	return out, nil
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
	diskPath := d.config.DataFolder(pathString)

	// Проверяем существование файла на диске
	if _, err := os.Stat(diskPath); err != nil {
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

func (d *diskService) Exists(path string) bool {
	path = normalizePath(path)
	if path == "" {
		return false
	}

	fsys, fullPath := d.resolvePath(path)
	if fsys != nil {
		if _, err := fs.Stat(fsys, fullPath); err == nil {
			return true
		} else if !errors.Is(err, fs.ErrNotExist) {
			return false
		}
	}

	diskPath := d.config.DataFolder(path)
	if _, err := os.Stat(diskPath); err == nil {
		return true
	}

	return false
}

func (d *diskService) GetTemplatesFS() fs.FS {
	return d.templatesFS
}

func (d *diskService) GetStaticFS() fs.FS {
	return d.staticFS
}

func (d *diskService) SetupStaticFiles(router *gin.Engine) error {
	if router == nil {
		return nil
	}

	// Добавляем заголовки Cache-Control для статических файлов
	router.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		if p == "/favicon.ico" ||
			strings.HasPrefix(p, "/static/") ||
			strings.HasPrefix(p, "/uploads/") ||
			strings.HasPrefix(p, "/public/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	})

	// /uploads/* - всегда с диска (пользовательские загруженные файлы)
	// Физически файлы в data/uploads, но URL остаются /uploads/
	router.Static("/uploads", d.config.DataFolder("uploads"))
	router.Static("/public", d.config.DataFolder("public"))

	// /static/* - только из embed (статические файлы)
	staticSub, err := fs.Sub(d.staticFS, "static")
	if err != nil {
		logger.Fatalf("[disk][SetupStaticFiles] fs.Sub(static) failed: %v", err)
		return err
	}
	router.StaticFS("/static", http.FS(staticSub))

	// favicon.ico из embed
	router.GET("/favicon.ico", func(c *gin.Context) {
		b, err := fs.ReadFile(d.staticFS, "static/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "image/x-icon", b)
	})

	return nil
}
