package service

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/file/models"
	"github.com/google/uuid"
)

type UploadService struct {
	fileService    *FileService
	storageService contract.Storage
}

func NewUploadService(
	fileService *FileService,
	storageService contract.Storage,
) *UploadService {
	return &UploadService{
		fileService,
		storageService,
	}
}

func (s *UploadService) SaveUploadedFile(file *multipart.FileHeader, folder string, user contract.User) (path string, err error) {
	ext, contentType := s.safeExt(file)
	if !s.isImage(contentType) {
		return "", fmt.Errorf("[file][UploadService][SaveUploadedFile] file:%s is not an image", file.Filename)
	}

	newUUID := uuid.New()
	if path, err = s.storageService.Save(file, folder, newUUID.String()+ext); err != nil {
		return
	}

	newFile := &models.File{
		UUID:         newUUID,
		UserID:       user.GetID(),
		File:         path,
		OriginalName: file.Filename,
		Size:         file.Size,
		Type:         contentType,
	}

	err = s.fileService.Create(newFile)
	if err != nil {
		newErr := s.storageService.Destroy(path)
		if newErr != nil {
			return "", newErr
		}

		return "", err
	}

	return
}

func (s *UploadService) SaveUploadedFiles(files []*multipart.FileHeader, dist string, user contract.User) ([]string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var paths []string
	newErr := errutil.New()

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader, dist string, user contract.User) {
			defer wg.Done()
			path, e := s.SaveUploadedFile(file, dist, user)
			if e != nil {
				newErr.Add(e)
				return
			}

			mu.Lock()
			paths = append(paths, path)
			mu.Unlock()
		}(file, dist, user)
	}

	wg.Wait()
	return paths, newErr.Error()
}

func (s *UploadService) DestroyFile(file string) error {
	return s.storageService.Destroy(file)
}

func (s *UploadService) Exist(file string) bool {
	return s.storageService.Exists(file)
}

func (s *UploadService) isImage(contentType string) bool {
	// Проверка MIME-типа на соответствие изображениям
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp":
		return true
	default:
		return false
	}
}

func (s *UploadService) safeExt(fh *multipart.FileHeader) (ext, contentType string) {
	// 1) из имени файла
	ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))
	if ext != "" && isAllowedExt(ext) {
		return "." + ext, mime.TypeByExtension("." + ext)
	}

	// 2) по заголовку из формы (если передавали)
	if ct := fh.Header.Get("Content-Type"); ct != "" {
		if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
			e := strings.ToLower(exts[0])
			if isAllowedExt(strings.TrimPrefix(e, ".")) {
				return e, ct
			}
		}
	}

	// 3) «понюхаем» первые 512 байт
	f, err := fh.Open()
	if err == nil {
		defer f.Close()
		var head [512]byte
		n, _ := io.ReadFull(f, head[:])
		ct := http.DetectContentType(head[:n])

		if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
			e := strings.ToLower(exts[0])
			if isAllowedExt(strings.TrimPrefix(e, ".")) {
				return e, ct
			}
		}

		// fallback на популярные типы
		switch {
		case strings.Contains(ct, "jpeg"):
			return ".jpg", ct
		case strings.Contains(ct, "png"):
			return ".png", ct
		case strings.Contains(ct, "gif"):
			return ".gif", ct
		case strings.Contains(ct, "webp"):
			return ".webp", ct
		case strings.Contains(ct, "pdf"):
			return ".pdf", ct
		}
	}

	// финальный дефолт
	return ".bin", "application/octet-stream"
}

func isAllowedExt(ext string) bool {
	switch ext {
	case "jpg", "jpeg", "png", "gif", "webp", "svg", "pdf", "mp4", "txt", "csv", "json":
		return true
	default:
		return false
	}
}
