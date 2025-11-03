package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
)

type InfoBlockProvider interface {
	GetForResourceUUID(resourceUUID string) []contract.InfoBlock
	DetachResourceUUID(resourceUUID string) error
	GetAll() []contract.InfoBlock
	Attach(infoBlockID uint, resourceUUID string) (infoBlocks []contract.InfoBlock, err error)
	SaveForm(block any, resourceUUID string) (contract.InfoBlock, error)
	SaveFormBatch(blocks []any, resourceUUID string) (infoBlock []contract.InfoBlock, err error)
	FindByTitle(title string) (contract.InfoBlock, error)
}
