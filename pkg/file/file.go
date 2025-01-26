package file

import (
	"errors"
	"fmt"
	"github.com/axlle-com/blog/pkg/app/config"
	"github.com/axlle-com/blog/pkg/app/logger"
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

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SaveUploadedFile(file *multipart.FileHeader, dist string) (path string, err error) {
	if !s.isImage(file) {
		return "", errors.New(fmt.Sprintf("Файл:%s не является изображением", file.Filename))
	}
	if path, err = s.save(file, dist); err != nil {
		return
	}
	return
}

func (s *Service) SaveUploadedFiles(files []*multipart.FileHeader, dist string) (paths []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader, dist string) {
			defer wg.Done()
			path, e := s.SaveUploadedFile(file, dist)
			if e != nil {
				logger.Error(e)
				return
			}

			mu.Lock()
			paths = append(paths, path)
			mu.Unlock()
		}(file, dist)
	}

	wg.Wait()
	return
}

func (s *Service) DeleteFile(file string) error {
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

func (s *Service) save(file *multipart.FileHeader, dst string) (string, error) {
	name := s.newName(dst, filepath.Ext(file.Filename))
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

func (s *Service) newName(dist, ext string) string {
	return fmt.Sprintf(config.Config().UploadPath()+"%s/%s%s", dist, uuid.New().String(), ext)
}

func (s *Service) realPath(path string) string {
	return config.Config().SrcFolderBuilder(path)
}

func (s *Service) isImage(fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
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
		return false
	}

	// Получение MIME-типа файла
	contentType := http.DetectContentType(buffer)

	// Проверка MIME-типа на соответствие изображениям
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp":
		return true
	default:
		return false
	}
}
