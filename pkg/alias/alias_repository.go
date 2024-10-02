package alias

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"gorm.io/gorm"
)

type Repository interface {
	GetByAlias(table, alias string, id uint) error
}

type repository struct {
	db *gorm.DB
}

func Repo() Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) GetByAlias(table, alias string, id uint) error {
	result := map[string]interface{}{}
	return r.db.Table(table).Where("alias = ?", alias).Where("id <> ?", id).Take(&result).Error
}
