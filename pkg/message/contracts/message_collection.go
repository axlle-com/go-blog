package contracts

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
)

type MessageCollectionService interface {
	GetAll() ([]*models.Message, error)
	GetAllIds() ([]uint, error)
	GetByIDs(ids []uint) ([]*models.Message, error)
	CountByField(field string, value any) (int64, error)
	Delete(messages []*models.Message) (err error)
	WithPaginate(p contracts.Paginator, filter *models.MessageFilter) ([]*models.Message, error)
	Paginator(p contracts.Paginator, filter *models.MessageFilter) (contracts.Paginator, error)
	Aggregates(messages []*models.Message) []*models.Message
}
