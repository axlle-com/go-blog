package contracts

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	WithTx(tx *gorm.DB) MessageRepository
	Create(message *models.Message) error
	GetByID(id uint) (*models.Message, error)
	GetByIDs(ids []uint) ([]*models.Message, error)
	CountByField(field string, value any) (int64, error)
	Update(message *models.Message) error
	Delete(*models.Message) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Message, error)
	GetAllIds() ([]uint, error)
	WithPaginate(p contracts.Paginator, filter *models.MessageFilter) ([]*models.Message, error)
	Paginator(paginator contracts.Paginator, filter *models.MessageFilter) (contracts.Paginator, error)
}
