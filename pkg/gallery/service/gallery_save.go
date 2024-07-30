package service

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
)

func SaveFromForm(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		return
	}
	form := c.Request.MultipartForm
	images := make(map[string]map[string]*models.GalleryImage)
	fileHeaders := make(map[string]*multipart.FileHeader)
	re := regexp.MustCompile(`^galleries\[(.+?)\]\[images\]\[(.+?)\]\[(.+)\]$`)
	for key, values := range form.Value {
		if matches := re.FindStringSubmatch(key); matches != nil {
			galleryID := matches[1]
			imageID := matches[2]
			field := matches[3]

			if _, ok := images[galleryID]; !ok {
				images[galleryID] = make(map[string]*models.GalleryImage)
			}

			if _, ok := images[galleryID][imageID]; !ok {
				images[galleryID][imageID] = &models.GalleryImage{}
			}

			image := images[galleryID][imageID]
			switch field {
			case "title":
				image.Title = &values[0]
			case "description":
				image.Description = &values[0]
			case "sort":
				image.Sort, _ = strconv.Atoi(values[0])
			}

		}
	}
	for _, headers := range form.File {
		for _, header := range headers {
			contentDisposition := header.Header.Get("Content-Disposition")
			re := regexp.MustCompile(`name="(galleries\[(.+?)\]\[images\]\[(.+?)\]\[file\])"`)
			if matches := re.FindStringSubmatch(contentDisposition); matches != nil {
				galleryID := matches[2]
				imageID := matches[3]

				if _, ok := images[galleryID]; !ok {
					images[galleryID] = make(map[string]*models.GalleryImage)
				}

				if _, ok := images[galleryID][imageID]; !ok {
					image := &models.GalleryImage{}
					images[galleryID][imageID] = image
				}

				if _, ok := fileHeaders[imageID]; !ok {
					fileHeaders[imageID] = header
				}
			}
		}
	}

	for galleryID, imagesMap := range images {
		galleryRepo := repository.NewGalleryRepository()
		var gallery *models.Gallery

		id, err := strconv.Atoi(galleryID)
		if err != nil {
			gallery = &models.Gallery{}
			err = galleryRepo.Create(gallery)
			if err != nil {
				continue
			}
		} else {
			gallery, err = galleryRepo.GetByID(uint(id))
			if err != nil {
				gallery = &models.Gallery{}
				err = galleryRepo.Create(gallery)
				if err != nil {
					continue
				}
			}
		}

		for imageID, image := range imagesMap {
			var imageOld *models.GalleryImage
			image.GalleryID = gallery.ID
			imageRepo := repository.NewGalleryImageRepository()
			id, err = strconv.Atoi(imageID)
			if err == nil {
				imageOld, err = imageRepo.GetByID(uint(id))
			}

			if header, ok := fileHeaders[imageID]; ok {
				newFileName := fmt.Sprintf("/public/uploads/%d/%s%s", gallery.ID, uuid.New().String(), filepath.Ext(header.Filename))
				if err := c.SaveUploadedFile(header, "src"+newFileName); err != nil {
					log.Println(err)
				}
				image.File = newFileName
				image.OriginalName = header.Filename
			}
			if imageOld == nil || imageOld.ID == 0 {
				err := imageRepo.Create(image)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				image.ID = imageOld.ID
				err := imageRepo.Update(image)
				if err != nil {
					log.Fatalln(err)
				}
			}
			log.Println("image: ", image)
		}
	}
}
