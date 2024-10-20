package file

import (
	"errors"
	"fmt"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/logger"
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

func SaveUploadedFile(file *multipart.FileHeader, dist string) (path string, err error) {
	if !isImageFile(file) {
		return "", errors.New(fmt.Sprintf("Файл:%s не является изображением", file.Filename))
	}
	path = fmt.Sprintf(config.Config().UploadPath()+"%s/%s%s", dist, uuid.New().String(), filepath.Ext(file.Filename))
	if err = save(file, realPath(path)); err != nil {
		return
	}
	return
}

func SaveUploadedFiles(files []*multipart.FileHeader, dist string) (uploadedFiles []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader, dist string) {
			defer wg.Done()
			path, e := SaveUploadedFile(file, dist)
			if e != nil {
				logger.Error(e)
				return
			}

			mu.Lock()
			uploadedFiles = append(uploadedFiles, path)
			mu.Unlock()
		}(file, dist)
	}

	wg.Wait()
	return
}

func DeleteFile(file string) error {
	if strings.HasPrefix(file, staticPath) {
		return nil
	}
	absPath, err := filepath.Abs(realPath(file))
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(absPath)
}

func save(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	_, err = io.Copy(out, src)
	return err
}

func realPath(path string) string {
	return config.Config().SrcFolderBuilder(path)
}

func isImageFile(fileHeader *multipart.FileHeader) bool {
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
