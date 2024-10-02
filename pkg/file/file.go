package file

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const staticPath = "/public/uploads/img/"

func SaveUploadedFile(file *multipart.FileHeader, dst string) (path string, err error) {
	cnf := config.GetConfig()
	path = fmt.Sprintf(cnf.UploadsPath+"%s/%s%s", dst, uuid.New().String(), filepath.Ext(file.Filename))
	if err := save(file, cnf.UploadsFolder+path); err != nil {
		return "", err
	}
	return path, nil
}

func DeleteFile(file string) error {
	cnf := config.GetConfig()
	if strings.HasPrefix(file, staticPath) {
		return nil
	}
	absPath, err := filepath.Abs(cnf.UploadsFolder + file)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return err
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
