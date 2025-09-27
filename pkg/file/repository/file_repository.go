package repository

import (
	"time"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/file/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	WithTx(tx *gorm.DB) FileRepository
	Create(file *models.File) error
	GetByID(id uint) (*models.File, error)
	GetByParams(params map[string]any, includeDeleted bool) ([]*models.File, error)
	GetByUUID(uuid uuid.UUID) (*models.File, error)
	GetByFile(string) (*models.File, error)
	GetByIDs(ids []uint) ([]*models.File, error)
	GetByUUIDs(uuids []uuid.UUID) ([]*models.File, error)
	Update(file *models.File) error
	Received(file []string) error
	Delete(id uint) error
	Destroy(id uint) error
	GetAll() ([]*models.File, error)
	GetAllIds() ([]uint, error)
}

type repository struct {
	db *gorm.DB
	*app.Paginate
}

func NewFileRepo(db *gorm.DB) FileRepository {
	return &repository{db: db}
}

func (r *repository) WithTx(tx *gorm.DB) FileRepository {
	newR := &repository{db: tx}
	return newR
}

func (r *repository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *repository) GetByID(id uint) (*models.File, error) {
	var file models.File
	if err := r.db.Select(file.Fields()).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *repository) GetByUUID(uuid uuid.UUID) (*models.File, error) {
	var file models.File
	if err := r.db.Select(file.Fields()).Where("uuid = ?", uuid).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *repository) GetByFile(path string) (*models.File, error) {
	var file models.File
	if err := r.db.Select(file.Fields()).Where("file = ?", path).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *repository) GetByIDs(ids []uint) ([]*models.File, error) {
	var files []*models.File
	if err := r.db.Select((&models.File{}).Fields()).Where("id IN (?)", ids).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *repository) GetByUUIDs(uuids []uuid.UUID) ([]*models.File, error) {
	var files []*models.File
	if err := r.db.Select((&models.File{}).Fields()).Where("uuid IN (?)", uuids).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *repository) GetByParams(params map[string]any, includeDeleted bool) ([]*models.File, error) {
	var files []*models.File
	db := r.db

	if includeDeleted {
		db = r.db.Unscoped()
	}

	if err := db.Where(params).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

func (r *repository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

func (r *repository) Received(files []string) error {
	if err := r.db.
		Model(&models.File{}).
		Where("file IN ?", files).
		Update("received_at", time.Now()).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

func (r *repository) Destroy(id uint) error {
	return r.db.Unscoped().Delete(&models.File{}, id).Error
}

func (r *repository) GetAll() ([]*models.File, error) {
	var files []*models.File
	if err := r.db.Select((&models.File{}).Fields()).Order("id ASC").Find(&files).Error; err != nil {
		return files, err
	}
	return files, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.File{}).Pluck("id", &ids).Error; err != nil {
		return ids, err
	}
	return ids, nil
}
