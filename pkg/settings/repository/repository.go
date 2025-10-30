package repository

import (
	"github.com/axlle-com/blog/pkg/settings/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	WithTx(tx *gorm.DB) Repository
	Get(namespace, key, scope string) (*models.Setting, error)
	Upsert(s *models.Setting) error
	ListNamespace(namespace, scope string) ([]models.Setting, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) WithTx(tx *gorm.DB) Repository { return &repo{db: tx} }

func (r *repo) Get(namespace, key, scope string) (*models.Setting, error) {
	var s models.Setting
	if err := r.db.
		Where("namespace = ? AND key = ? AND scope = ?", namespace, key, scope).
		First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repo) Upsert(s *models.Setting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "namespace"}, {Name: "key"}, {Name: "scope"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "value", "owner_uuid", "sort", "updated_at"}),
	}).Create(s).Error
}

func (r *repo) ListNamespace(namespace, scope string) ([]models.Setting, error) {
	var list []models.Setting
	q := r.db.Where("namespace = ?", namespace)
	if scope != "" {
		q = q.Where("scope = ?", scope)
	}
	if err := q.Order("sort ASC, created_at ASC, key ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
