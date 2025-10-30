package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type GalleryProvider interface {
	GetForResourceUUID(resourceUUID string) []contract.Gallery
	GetIndexesForResources(resources []contract.Resource) map[uuid.UUID][]contract.Gallery
	GetAll() []contract.Gallery
	SaveForm(g any, resource contract.Resource) (contract.Gallery, error)
	SaveFormBatch(anys []any, resource contract.Resource) (galleries []contract.Gallery, err error)
	DetachResource(contract.Resource) error
}

type ImageProvider interface {
	GetForGallery(id uint) []contract.Image
	GetAll() []contract.Image
}
