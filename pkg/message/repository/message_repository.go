package repository

import (
	"fmt"
	"github.com/axlle-com/blog/app/db"
	app "github.com/axlle-com/blog/app/models"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	"gorm.io/gorm"
	"log"
)

type MessageRepository interface {
	WithTx(tx *gorm.DB) MessageRepository
	Create(message *models.Message) error
	GetByID(id uint) (*models.Message, error)
	GetByIDs(ids []uint) ([]*models.Message, error)
	Update(message *models.Message) error
	Delete(*models.Message) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Message, error)
	GetAllIds() ([]uint, error)
	WithPaginate(p contracts2.Paginator, filter *models.MessageFilter) ([]*models.Message, error)
}

type repository struct {
	db *gorm.DB
	*app.Paginate
}

func NewMessageRepo() MessageRepository {
	r := &repository{db: db.GetDB()}
	return r
}

func (r *repository) WithTx(tx *gorm.DB) MessageRepository {
	newR := &repository{db: tx}
	return newR
}

func (r *repository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *repository) GetByID(id uint) (*models.Message, error) {
	var message models.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *repository) GetByIDs(ids []uint) ([]*models.Message, error) {
	var messages []*models.Message
	if err := r.db.Where("id IN (?)", ids).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *repository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *repository) Delete(message *models.Message) error {
	return r.db.Delete(&models.Message{}, message.ID).Error
}

func (r *repository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.Message{}).Error
}

func (r *repository) GetAll() ([]*models.Message, error) {
	var message []*models.Message
	if err := r.db.Order("id ASC").Find(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Message{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("Failed to fetch IDs from the database: %v\n", err)
	}
	return ids, nil
}

func (r *repository) WithPaginate(p contracts2.Paginator, filter *models.MessageFilter) ([]*models.Message, error) {
	var messages []*models.Message
	var total int64

	message := models.Message{}
	table := message.GetTable()

	query := r.db.Model(&message)

	// TODO WHERE IN; LIKE
	for col, val := range filter.GetMap() {
		if col == "title" {
			query = query.Where(fmt.Sprintf("%s.%v ilike ?", table, col), fmt.Sprintf("%%%v%%", val))
			continue
		}
		query = query.Where(fmt.Sprintf("%s.%v = ?", table, col), val)
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order(fmt.Sprintf("%s.id ASC", message.GetTable())).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))
	return messages, nil
}
