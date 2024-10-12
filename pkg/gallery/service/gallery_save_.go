package service

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type GalleriesCollections map[string]ImagesCollections
type ImagesCollections map[string]*models.Image

func _SaveFromForm(c *gin.Context) []*models.Gallery {
	err := c.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		logger.Error(err)
		return nil
	}

	title, _ := c.Get("title")
	titleStr, _ := title.(string)

	images := parseFormValue(c)
	images = parseFormFile(c, images)

	var group sync.WaitGroup
	var mutex sync.Mutex
	var galleries []*models.Gallery
	for galleryID, imagesMap := range images {
		galleryRepo := models.GalleryRepo()
		var gallery = &models.Gallery{Title: &titleStr}
		id, err := strconv.Atoi(galleryID)
		if err != nil {
			err = galleryRepo.Create(gallery)
			if err != nil {
				logger.Error(err)
				continue
			}
		} else {
			gallery, err = galleryRepo.GetByID(uint(id))
			if err != nil {
				err = galleryRepo.Create(gallery)
				if err != nil {
					logger.Error(err)
					continue
				}
			}
		}
		galleries = append(galleries, gallery)

		for i, image := range imagesMap {
			image.GalleryID = gallery.ID
			id, _ = strconv.Atoi(i)
			image.ID = uint(id)
			group.Add(1)
			image := image
			go func() {
				defer group.Done()
				_, err := SaveImage(image)
				if err != nil {
					logger.Error(err)
				} else {
					mutex.Lock()
					gallery.Images = append(gallery.Images, image)
					mutex.Unlock()
				}
			}()
		}
	}
	group.Wait()
	return galleries
}

func parseFormValue(c *gin.Context) GalleriesCollections {
	form := c.Request.MultipartForm
	re := regexp.MustCompile(`^galleries\[(.+?)\]\[images\]\[(.+?)\]\[(.+)\]$`)
	images := make(GalleriesCollections)

	for key, values := range form.Value {
		if matches := re.FindStringSubmatch(key); matches != nil {
			galleryID := matches[1]
			imageID := matches[2]
			field := matches[3]

			if _, ok := images[galleryID]; !ok {
				images[galleryID] = make(ImagesCollections)
			}

			if _, ok := images[galleryID][imageID]; !ok {
				images[galleryID][imageID] = &models.Image{}
			}

			image := images[galleryID][imageID]
			switch field {
			case "title":
				image.Title = db.StrPtr(values[0])
			case "description":
				image.Description = db.StrPtr(values[0])
			case "sort":
				image.Sort, _ = strconv.Atoi(values[0])
			}
		}
	}
	return images
}

func parseFormFile(c *gin.Context, images GalleriesCollections) GalleriesCollections {
	form := c.Request.MultipartForm
	for _, headers := range form.File {
		for _, header := range headers {
			if !isImageFile(header) {
				continue
			}
			contentDisposition := header.Header.Get("Content-Disposition")
			re := regexp.MustCompile(`name="(galleries\[(.+?)\]\[images\]\[(.+?)\]\[file\])"`)
			if matches := re.FindStringSubmatch(contentDisposition); matches != nil {
				galleryID := matches[2]
				imageID := matches[3]

				if _, ok := images[galleryID]; !ok {
					images[galleryID] = make(ImagesCollections)
				}

				if _, ok := images[galleryID][imageID]; !ok {
					image := &models.Image{FileHeader: header}
					images[galleryID][imageID] = image
				} else {
					images[galleryID][imageID].FileHeader = header
				}
			}
		}
	}
	return images
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
