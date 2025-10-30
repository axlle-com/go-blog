package provider

import (
	"sync"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/google/uuid"
)

func NewProvider(
	gallery repository.GalleryRepository,
	service *service.GalleryService,
) apppPovider.GalleryProvider {
	return &provider{
		gallery: gallery,
		service: service,
	}
}

type provider struct {
	gallery repository.GalleryRepository
	service *service.GalleryService
}

func (p *provider) GetForResourceUUID(resourceUUID string) []contract.Gallery {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return nil
	}

	galleries, err := p.gallery.WithImages().GetForResource(newUUID)
	collection := make([]contract.Gallery, 0, len(galleries))
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Errorf("[gallery][provider][GetForResourceUUID] error: %v", err)
	return nil
}

func (p *provider) GetIndexesForResources(resources []contract.Resource) map[uuid.UUID][]contract.Gallery {
	uuids := make([]uuid.UUID, 0, len(resources))
	for _, resource := range resources {
		uuids = append(uuids, resource.GetUUID())
	}

	galleries, err := p.gallery.WithImages().GetForResources(uuids)
	collection := make(map[uuid.UUID][]contract.Gallery)
	if err == nil {
		for _, gallery := range galleries {
			if _, ok := collection[gallery.GetResourceUUID()]; !ok {
				collection[gallery.GetResourceUUID()] = make([]contract.Gallery, 0)
			}
			collection[gallery.GetResourceUUID()] = append(collection[gallery.GetResourceUUID()], gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) DetachResource(resource contract.Resource) (err error) {
	err = p.service.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contract.Gallery {
	var collection []contract.Gallery
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

func (p *provider) SaveForm(g any, resource contract.Resource) (gallery contract.Gallery, err error) {
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

func (p *provider) SaveFormBatch(anys []any, resource contract.Resource) (galleries []contract.Gallery, err error) {
	var wg sync.WaitGroup

	// Блокировки для конкурентного доступа к срезам.
	var galleriesMu sync.Mutex
	var errorsMu sync.Mutex // @todo new error

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

	newErr := errutil.New()
	if len(errs) > 0 {
		for _, err := range errs {
			newErr.Add(err)
		}
	}

	return galleries, newErr.Error()
}
