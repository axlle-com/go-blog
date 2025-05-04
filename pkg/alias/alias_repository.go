package alias

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"gorm.io/gorm"
)

type AliasRepository interface {
	GetByAlias(id uint, table, alias string) error
	WithTx(tx *gorm.DB) AliasRepository
}

type repository struct {
	db *gorm.DB
}

func NewAliasRepo(db contracts.DB) AliasRepository {
	return &repository{db: db.GORM()}
}

func (r *repository) WithTx(tx *gorm.DB) AliasRepository {
	return &repository{db: tx}
}

func (r *repository) GetByAlias(id uint, table, alias string) error {
	result := map[string]interface{}{}
	return r.db.Table(table).Where("alias = ?", alias).Where("id <> ?", id).Take(&result).Error
}
