package service

import (
	"errors"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

type ImageService struct {
	imageRepo  repository.GalleryImageRepository
	imageEvent *ImageEvent
	api        *api.Api
}

func NewImageService(
	image repository.GalleryImageRepository,
	imageEvent *ImageEvent,
	api *api.Api,
) *ImageService {
	return &ImageService{
		imageRepo:  image,
		imageEvent: imageEvent,
		api:        api,
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

func (s *ImageService) DeleteImages(images []*models.Image) (err error) {
	var ids []uint
	var resImages []*models.Image
	errCollection := errutil.New()
	for _, image := range images {
		if err := s.imageEvent.DeletingImage(image); err == nil {
			ids = append(ids, image.ID)
			resImages = append(resImages, image)
		} else {
			errCollection.Add(err)
		}
	}

	if err = errCollection.Error(); err != nil {
		return err
	}

	if len(ids) > 0 {
		if err = s.imageRepo.DeleteByIDs(ids); err == nil {
			for _, image := range resImages {
				if err := s.imageEvent.DeletedImage(image); err != nil {
					errCollection.Add(err)
				}
			}

			if err = errCollection.Error(); err != nil {
				return errors.New("image deletion errors occurred")
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
