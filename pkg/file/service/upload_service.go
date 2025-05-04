package service

import (
	"fmt"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/file/models"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const staticPath = "/public/img/"

type UploadService struct {
	service *Service
}

func NewUploadService(service *Service) *UploadService {
	return &UploadService{service}
}

func (s *UploadService) SaveUploadedFile(file *multipart.FileHeader, dist string, user contracts.User) (path string, err error) {
	contentType := s.contentType(file)
	if !s.isImage(contentType) {
		return "", fmt.Errorf("Файл:%s не является изображением", file.Filename)
	}

	newUUID := uuid.New()
	if path, err = s.save(file, dist, newUUID.String()); err != nil {
		return
	}

	newFile := &models.File{
		UUID:         newUUID,
		UserID:       user.GetID(),
		File:         path,
		OriginalName: file.Filename,
		Size:         file.Size,
		Type:         contentType,
		IsReceived:   false,
	}

	err = s.service.Create(newFile)
	if err != nil {
		newErr := s.DestroyFile(path)
		if newErr != nil {
			return "", newErr
		}
		return "", err
	}

	return
}

func (s *UploadService) SaveUploadedFiles(files []*multipart.FileHeader, dist string, user contracts.User) (paths []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader, dist string, user contracts.User) {
			defer wg.Done()
			path, e := s.SaveUploadedFile(file, dist, user)
			if e != nil {
				logger.Error(e)
				return
			}

			mu.Lock()
			paths = append(paths, path)
			mu.Unlock()
		}(file, dist, user)
	}

	wg.Wait()
	return
}

func (s *UploadService) DestroyFile(file string) error {
	if strings.HasPrefix(file, staticPath) {
		return nil
	}
	absPath, err := filepath.Abs(s.realPath(file))
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(absPath)
}

func (s *UploadService) Exist(file string) bool {
	if strings.HasPrefix(file, staticPath) {
		return true
	}

	absPath, err := filepath.Abs(s.realPath(file))
	if err != nil {
		// на всякий случай, если не получилось построить путь — говорим, что нет
		return false
	}

	if info, err := os.Stat(absPath); err == nil {
		// Убедимся, что это не директория (если нужно только файлы)
		return !info.IsDir()
	} else if os.IsNotExist(err) {
		// Файл точно отсутствует
		return false
	} else {
		// Какая-то другая ошибка (например, прав доступа) — можно трактовать как “нет”
		return false
	}
}

func (s *UploadService) save(file *multipart.FileHeader, dst, newUUID string) (string, error) {
	name := s.newName(dst, filepath.Ext(file.Filename), newUUID)
	path := s.realPath(name)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	if err = os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", err
	}

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			logger.Error(err)
		}
	}(out)

	_, err = io.Copy(out, src)
	return name, err
}

func (s *UploadService) newName(dist, ext, newUUID string) string {
	return fmt.Sprintf(config.Config().UploadPath()+"%s/%s%s", dist, newUUID, ext)
}

func (s *UploadService) realPath(path string) string {
	return config.Config().SrcFolderBuilder(path)
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

func (s *UploadService) contentType(fileHeader *multipart.FileHeader) string {
	file, err := fileHeader.Open()
	if err != nil {
		return ""
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logger.Error(err)
		}
	}(file)

	// Чтение первых 512 байт файла для определения MIME-типа
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return ""
	}

	// Получение MIME-типа файла
	return http.DetectContentType(buffer)
}
