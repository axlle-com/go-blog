package repository

import (
	"errors"
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"gorm.io/gorm"
)

type InfoBlockRepository interface {
	WithTx(tx *gorm.DB) InfoBlockRepository
	Create(infoBlock *models.InfoBlock) error
	GetByID(id uint) (*models.InfoBlock, error)
	Update(infoBlock *models.InfoBlock) error
	Delete(infoBlock *models.InfoBlock) error
	GetAll() ([]*models.InfoBlock, error)
	GetByIDs(ids []uint) ([]*models.InfoBlock, error)
	DeleteByIDs(ids []uint) (err error)
	GetForResource(resource contracts.Resource) ([]*models.InfoBlock, error)
}

type infoBlockRepository struct {
	db *gorm.DB
}

func NewInfoBlockRepo() InfoBlockRepository {
	r := &infoBlockRepository{db: db.GetDB()}
	return r
}

func (r *infoBlockRepository) WithTx(tx *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: tx}
}

func (r *infoBlockRepository) Create(infoBlock *models.InfoBlock) error {
	infoBlock.Creating()
	return r.db.Create(infoBlock).Error
}

func (r *infoBlockRepository) GetByID(id uint) (*models.InfoBlock, error) {
	var model models.InfoBlock
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *infoBlockRepository) Update(infoBlock *models.InfoBlock) error {
	infoBlock.Updating()
	return r.db.Select(
		"TemplateID",
		"MetaTitle",
		"Media",
		"Title",
		"Description",
		"Image",
	).Save(infoBlock).Error
}

func (r *infoBlockRepository) Delete(infoBlock *models.InfoBlock) error {
	if infoBlock.Deleting() {
		return r.db.Delete(&models.InfoBlock{}, infoBlock.ID).Error
	}
	return errors.New("при удалении произошли ошибки")
}

func (r *infoBlockRepository) GetAll() ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	if err := r.db.Find(&infoBlocks).Error; err != nil {
		return nil, err
	}
	return infoBlocks, nil
}

func (r *infoBlockRepository) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	query := r.db.Where("id IN ?", ids)

	if err := query.Find(&infoBlocks).Error; err != nil {
		return nil, err
	}
	return infoBlocks, nil
}

func (r *infoBlockRepository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.InfoBlock{}).Error
}

func (r *infoBlockRepository) GetForResource(resource contracts.Resource) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	query := r.db.
		Joins("inner join info_block_has_resources as r on info_blocks.id = r.info_block_id").
		Where("r.resource_uuid = ?", resource.GetUUID()).
		Model(&models.InfoBlock{})

	err := query.Find(&infoBlocks).Error
	return infoBlocks, err
}
