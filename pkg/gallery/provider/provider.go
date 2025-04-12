package provider

import (
	"errors"
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type GalleryProvider interface {
	GetForResource(contracts2.Resource) []contracts2.Gallery
	GetForResources([]contracts2.Resource) []contracts2.Gallery
	GetIndexesForResources(resources []contracts2.Resource) map[uuid.UUID][]contracts2.Gallery
	GetAll() []contracts2.Gallery
	SaveForm(g any, resource contracts2.Resource) (contracts2.Gallery, error)
	SaveFormBatch(anys []any, resource contracts2.Resource) (galleries []contracts2.Gallery, err error)
	DetachResource(contracts2.Resource) error
}

func NewProvider(
	gallery repository.GalleryRepository,
	service *service.GalleryService,
) GalleryProvider {
	return &provider{
		gallery: gallery,
		service: service,
	}
}

type provider struct {
	gallery repository.GalleryRepository
	service *service.GalleryService
}

func (p *provider) GetForResource(resource contracts2.Resource) []contracts2.Gallery {
	galleries, err := p.gallery.WithImages().GetForResource(resource.GetUUID())
	collection := make([]contracts2.Gallery, 0, len(galleries))
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetForResources(resources []contracts2.Resource) []contracts2.Gallery {
	uuids := make([]uuid.UUID, 0, len(resources))
	for _, resource := range resources {
		uuids = append(uuids, resource.GetUUID())
	}

	galleries, err := p.gallery.WithImages().GetForResources(uuids)
	collection := make([]contracts2.Gallery, 0, len(galleries))
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetIndexesForResources(resources []contracts2.Resource) map[uuid.UUID][]contracts2.Gallery {
	uuids := make([]uuid.UUID, 0, len(resources))
	for _, resource := range resources {
		uuids = append(uuids, resource.GetUUID())
	}

	galleries, err := p.gallery.WithImages().GetForResources(uuids)
	collection := make(map[uuid.UUID][]contracts2.Gallery)
	if err == nil {
		for _, gallery := range galleries {
			if _, ok := collection[gallery.GetResourceUUID()]; !ok {
				collection[gallery.GetResourceUUID()] = make([]contracts2.Gallery, 0)
			}
			collection[gallery.GetResourceUUID()] = append(collection[gallery.GetResourceUUID()], gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) DetachResource(resource contracts2.Resource) (err error) {
	err = p.service.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contracts2.Gallery {
	var collection []contracts2.Gallery
	galleries, err := p.gallery.GetAll()
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveForm(g any, resource contracts2.Resource) (gallery contracts2.Gallery, err error) {
	gal := app.LoadStruct(&models.Gallery{}, g).(*models.Gallery)
	if gal.ID == 0 {
		gallery, err = p.service.CreateGallery(gal)
	} else {
		gallery, err = p.service.UpdateGallery(gal)
	}

	if err != nil {
		return nil, err
	}

	err = p.service.Attach(resource, gallery)
	if err != nil {
		return nil, err
	}

	return gallery, nil
}

func (p *provider) SaveFormBatch(anys []any, resource contracts2.Resource) (galleries []contracts2.Gallery, err error) {
	var wg sync.WaitGroup

	// Блокировки для конкурентного доступа к срезам.
	var galleriesMu sync.Mutex
	var errorsMu sync.Mutex

	var errs []error

	for _, gall := range anys {
		wg.Add(1)
		// Передаем gall как параметр, чтобы избежать проблем замыкания.
		go func(g any) {
			defer wg.Done()
			var localErr error
			var gallery *models.Gallery
			gal := app.LoadStruct(&models.Gallery{}, g).(*models.Gallery)

			// Если галерея новая, создаем, иначе обновляем
			if gal.ID == 0 {
				gallery, localErr = p.service.CreateGallery(gal)
			} else {
				gallery, localErr = p.service.UpdateGallery(gal)
			}

			if localErr != nil {
				logger.Error(localErr)
				errorsMu.Lock()
				errs = append(errs, localErr)
				errorsMu.Unlock()
				return
			}

			localErr = p.service.Attach(resource, gallery)
			if localErr != nil {
				logger.Error(localErr)
				errorsMu.Lock()
				errs = append(errs, localErr)
				errorsMu.Unlock()
				return
			}

			// Если все успешно, добавляем галерею в итоговый срез.
			galleriesMu.Lock()
			galleries = append(galleries, gallery)
			galleriesMu.Unlock()
		}(gall)
	}

	wg.Wait()

	if len(errs) > 0 {
		errMsg := ""
		for _, e := range errs {
			errMsg += e.Error() + "; "
		}
		errMsg = strings.TrimSuffix(errMsg, "; ")
		err = errors.New(errMsg)
	}

	return galleries, err
}
