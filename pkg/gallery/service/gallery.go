package service

import (
	"errors"
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GalleryService struct {
	galleryRepo  repository.GalleryRepository
	galleryEvent *GalleryEvent
	imageService *ImageService
	resourceRepo repository.GalleryResourceRepository
	fileProvider fileProvider.FileProvider
}

func NewGalleryService(
	galleryRepo repository.GalleryRepository,
	galleryEvent *GalleryEvent,
	imageService *ImageService,
	resourceRepo repository.GalleryResourceRepository,
	fileProvider fileProvider.FileProvider,
) *GalleryService {
	return &GalleryService{
		galleryRepo:  galleryRepo,
		galleryEvent: galleryEvent,
		imageService: imageService,
		resourceRepo: resourceRepo,
		fileProvider: fileProvider,
	}
}

func (s *GalleryService) CreateGallery(gallery *models.Gallery) (*models.Gallery, error) {
	if err := s.galleryRepo.Create(gallery); err != nil {
		return nil, err
	}

	err := s.galleryImageUpdate(gallery)
	return gallery, err
}

func (s *GalleryService) UpdateGallery(gallery *models.Gallery) (*models.Gallery, error) {
	if err := s.galleryRepo.Update(gallery); err != nil {
		return nil, err
	}

	err := s.galleryImageUpdate(gallery)
	return gallery, err
}

func (s *GalleryService) Attach(resource contracts2.Resource, gallery contracts2.Gallery) error {
	hasRepo, err := s.resourceRepo.GetByParams(resource.GetUUID(), gallery.GetID())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.GalleryHasResource{
				ResourceUUID: resource.GetUUID(),
				GalleryID:    gallery.GetID(),
			},
		)
		return err
	}
	return nil
}

func (s *GalleryService) DeleteForResource(resource contracts2.Resource) error {
	byResource, err := s.resourceRepo.GetByResource(resource)
	if err != nil {
		return err
	}
	if len(byResource) == 0 {
		return nil
	}

	// Группируем записи по GalleryID
	galleryResources := make(map[uint][]uuid.UUID)
	for _, res := range byResource {
		galleryResources[res.GalleryID] = append(galleryResources[res.GalleryID], res.ResourceUUID)
	}

	var detachGalleryIDs []uint
	var deleteGalleryIDs []uint

	// Определяем для каждой галереи, сколько ресурсов ей принадлежит
	for galleryID, resources := range galleryResources {
		if len(resources) > 1 {
			detachGalleryIDs = append(detachGalleryIDs, galleryID)
		} else {
			deleteGalleryIDs = append(deleteGalleryIDs, galleryID)
		}
	}

	if len(detachGalleryIDs) > 0 {
		for _, galleryID := range detachGalleryIDs {
			err = s.resourceRepo.DeleteByParams(resource.GetUUID(), galleryID)
			if err != nil {
				return err
			}
		}
	}

	if len(deleteGalleryIDs) > 0 {
		galleries, err := s.galleryRepo.WithImages().GetByIDs(deleteGalleryIDs)
		if err != nil {
			return err
		}
		err = s.DeleteGalleries(galleries)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GalleryService) DeleteGalleries(galleries []*models.Gallery) (err error) {
	var ids []uint
	for _, gallery := range galleries {
		if err = s.galleryEvent.DeletingGallery(gallery); err != nil {
			return err
		}
		ids = append(ids, gallery.ID)
	}

	if len(ids) > 0 {
		if err = s.galleryRepo.DeleteByIDs(ids); err == nil {
			for _, gallery := range galleries {
				if err = s.galleryEvent.DeletedGallery(gallery); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return err
}

func (s *GalleryService) galleryImageUpdate(gallery *models.Gallery) error {
	if len(gallery.Images) == 0 {
		return nil
	}

	var err error
	var file string
	var errSlice []error
	slice := make([]*models.Image, 0)
	sliceFiles := make([]string, 0)

	for _, item := range gallery.Images {
		if item == nil {
			continue
		}

		if item.ID == 0 {
			file = item.File
		}

		item.GalleryID = gallery.ID
		image, err := s.imageService.SaveImage(item)
		if err != nil {
			logger.Errorf("[GalleryService][galleryImageUpdate] Error: %v", err)
			errSlice = append(errSlice, err)
			continue
		}
		if image == nil {
			continue
		}

		slice = append(slice, image)

		if file != "" {
			sliceFiles = append(sliceFiles, file)
		}
	}

	if len(sliceFiles) > 0 {
		err = s.fileProvider.Received(sliceFiles)
		if err != nil {
			logger.Errorf("[GalleryService][galleryImageUpdate] Error: %v", err)
			errSlice = append(errSlice, err)
		}
	}

	if len(errSlice) > 0 {
		err = errors.New("были ошибки при сохранении изображения")
	}

	gallery.Images = slice

	return err
}
