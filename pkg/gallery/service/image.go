package service

import (
	"errors"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

type ImageService struct {
	imageRepo  repository.GalleryImageRepository
	imageEvent *ImageEvent
}

func NewImageService(
	image repository.GalleryImageRepository,
	imageEvent *ImageEvent,
) *ImageService {
	return &ImageService{
		imageRepo:  image,
		imageEvent: imageEvent,
	}
}

func (s *ImageService) SaveImage(image *models.Image) (*models.Image, error) {
	if image.ID == 0 {
		if err := s.imageRepo.Create(image); err != nil {
			return nil, err
		}
	} else {
		if err := s.imageRepo.Update(image); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func (s *ImageService) DeleteImages(is []*models.Image) (err error) {
	var ids []uint
	var resImages []*models.Image
	var isErr bool
	for _, im := range is {
		if err := s.imageEvent.DeletingImage(im); err == nil {
			ids = append(ids, im.ID)
			resImages = append(resImages, im)
		} else {
			isErr = true
			logger.Error(err)
		}
	}

	if isErr {
		return errors.New("Ошибки при удалении изображений")
	}

	if len(ids) > 0 {
		if err = s.imageRepo.DeleteByIDs(ids); err == nil {
			for _, im := range resImages {
				if err := s.imageEvent.DeletedImage(im); err != nil {
					logger.Error(err)
					isErr = true
				}
			}

			if isErr {
				return errors.New("Ошибки при удалении изображений")
			}
		}
	}
	return
}

func (s *ImageService) DeleteImage(im *models.Image) (err error) {
	if err = s.imageEvent.DeletingImage(im); err != nil {
		return
	}

	if err = s.imageRepo.Delete(im); err == nil {
		if err = s.imageEvent.DeletedImage(im); err != nil {
			return
		}
		return
	}
	return
}
