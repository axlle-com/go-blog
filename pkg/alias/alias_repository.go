package alias

import (
	"github.com/axlle-com/blog/app/db"
	"gorm.io/gorm"
)

type AliasRepository interface {
	GetByAlias(id uint, table, alias string) error
	WithTx(tx *gorm.DB) AliasRepository
}

type repository struct {
	db *gorm.DB
}

func NewAliasRepo() AliasRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) WithTx(tx *gorm.DB) AliasRepository {
	return &repository{db: tx}
}

func (r *repository) GetByAlias(id uint, table, alias string) error {
	result := map[string]interface{}{}
	return r.db.Table(table).Where("alias = ?", alias).Where("id <> ?", id).Take(&result).Error
}
