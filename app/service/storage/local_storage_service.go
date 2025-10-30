package storage

import (
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/axlle-com/blog/app/models/contract"
)

const staticPath = "/public/img/"

type LocalStorageService struct {
	publicBase      string // https://site.com
	localPathPrefix string // web-путь, например: /public/uploads
	fsBase          string // абсолютный путь на диске, например: /abs/path/src/public/uploads
}

func NewLocalStorageService(cfg contract.Config) contract.Storage {
	fsBase := filepath.Join(cfg.SrcFolder(), cfg.UploadPath())

	return &LocalStorageService{
		publicBase:      strings.TrimRight(cfg.AppHost(), "/"),
		localPathPrefix: cfg.UploadPath(),
		fsBase:          fsBase,
	}
}

func (b *LocalStorageService) Save(fh *multipart.FileHeader, folder, fileName string) (string, error) {
	rel := joinURLPath(b.localPathPrefix, folder, fileName) // /public/uploads/posts/uuid.jpg
	fullURL := b.publicBase + rel

	absPath := filepath.Join(b.fsBase, folder) // <-- ПРАВИЛЬНЫЙ корень записи
	if err := os.MkdirAll(absPath, 0o750); err != nil {
		return "", err
	}

	dst := filepath.Join(absPath, fileName)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}
	return fullURL, nil
}

func (b *LocalStorageService) SaveReader(reader io.Reader, _ int64, folder, filename, _ string) (string, error) {
	rel := joinURLPath(b.localPathPrefix, folder, filename)
	fullURL := b.publicBase + rel

	absPath := filepath.Join(b.fsBase, folder) // <-- тоже из fsBase
	if err := os.MkdirAll(absPath, 0o750); err != nil {
		return "", err
	}
	dst := filepath.Join(absPath, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, reader); err != nil {
		return "", err
	}
	return fullURL, nil
}

func (b *LocalStorageService) Destroy(urlOrPath string) error {
	p := b.toAbsPath(urlOrPath)
	if p == "" {
		return nil
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(p)
}

func (b *LocalStorageService) Exists(urlOrPath string) bool {
	if isStaticRef(urlOrPath) {
		return true
	}
	p := b.toAbsPath(urlOrPath)
	if p == "" {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func (b *LocalStorageService) toAbsPath(urlOrPath string) string {
	// Полный URL → срежем localPathPrefix и прибьём к fsBase
	if u, err := url.Parse(urlOrPath); err == nil && u.Scheme != "" && u.Host != "" {
		rel := strings.TrimPrefix(u.Path, b.localPathPrefix)
		rel = strings.TrimLeft(rel, "/")
		if rel == "" {
			return ""
		}
		return filepath.Join(b.fsBase, rel)
	}

	// Относительный web-путь с префиксом (/public/uploads/...)
	if strings.HasPrefix(urlOrPath, b.localPathPrefix) {
		rel := strings.TrimPrefix(urlOrPath, b.localPathPrefix)
		rel = strings.TrimLeft(rel, "/")
		if rel == "" {
			return ""
		}
		return filepath.Join(b.fsBase, rel)
	}

	// Абсолютный путь в ФС оставляем как есть
	if filepath.IsAbs(urlOrPath) {
		return urlOrPath
	}

	return ""
}

func joinURLPath(parts ...string) string {
	var cleaned []string
	for _, p := range parts {
		cleaned = append(cleaned, strings.Trim(p, "/"))
	}
	return "/" + strings.Join(cleaned, "/")
}

func isStaticRef(s string) bool {
	if s == "" {
		return false
	}
	// 1) Полный URL c /public/img/ в path
	if u, err := url.Parse(s); err == nil && u.Scheme != "" && u.Host != "" {
		return strings.HasPrefix(u.Path, staticPath)
	}
	// 2) Относительный web-путь
	asSlash := filepath.ToSlash(s)
	if strings.HasPrefix(asSlash, staticPath) {
		return true
	}
	// 3) Абсолютный путь в ФС, содержащий /public/img/
	return strings.Contains(asSlash, staticPath)
}
