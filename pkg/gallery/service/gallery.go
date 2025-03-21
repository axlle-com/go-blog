package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"gorm.io/gorm"
)

type GalleryService struct {
	galleryRepo  repository.GalleryRepository
	galleryEvent *GalleryEvent
	imageService *ImageService
	resourceRepo repository.GalleryResourceRepository
}

func NewGalleryService(
	galleryRepo repository.GalleryRepository,
	galleryEvent *GalleryEvent,
	imageService *ImageService,
	resourceRepo repository.GalleryResourceRepository,
) *GalleryService {
	return &GalleryService{
		galleryRepo:  galleryRepo,
		galleryEvent: galleryEvent,
		imageService: imageService,
		resourceRepo: resourceRepo,
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

func (s *GalleryService) Attach(resource contracts.Resource, gallery contracts.Gallery) error {
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
	}
	return nil
}

func (s *GalleryService) DeleteForResource(resource contracts.Resource) (err error) {
	byResource, err := s.resourceRepo.GetByResource(resource)
	if err != nil {
		return err
	}

	all := make(map[uint]*models.GalleryHasResource)
	only := make(map[uint]*models.GalleryHasResource)
	detach := make(map[uint]*models.GalleryHasResource)
	var galleryIDs []uint
	if byResource == nil {
		return nil
	}

	for _, r := range byResource {
		if r.ResourceUUID != resource.GetUUID() {
			all[r.GalleryID] = r
		} else {
			only[r.GalleryID] = r
		}
	}

	for id, _ := range only {
		if _, ok := all[id]; ok {
			detach[id] = all[id]
		} else {
			galleryIDs = append(galleryIDs, id)
		}
	}

	if len(detach) > 0 { // TODO need test
		for _, r := range detach {
			err = s.resourceRepo.DeleteByParams(r.ResourceUUID, r.GalleryID)
			if err != nil {
				return err
			}
		}
	}

	if len(galleryIDs) > 0 {
		galleries, err := s.galleryRepo.WithImages().GetByIDs(galleryIDs)
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
	var err error
	if len(gallery.Images) > 0 {
		slice := make([]*models.Image, 0)
		var eSlice []error
		for _, item := range gallery.Images {
			if item == nil {
				continue
			}
			item.GalleryID = gallery.ID
			image, e := s.imageService.SaveImage(item)
			if e != nil {
				eSlice = append(eSlice, e)
				continue
			}
			if image == nil {
				continue
			}

			slice = append(slice, image)
		}
		if len(eSlice) > 0 {
			err = errors.New("были ошибки при сохранении изображения")
		}
		gallery.Images = slice
	}

	return err
}
