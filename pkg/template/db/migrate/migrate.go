package migrate

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/template/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator() contracts.Migrator {
	return &migrator{db: db.GetDB()}
}

func (m *migrator) Migrate() {
	err := m.db.AutoMigrate(
		&models.Template{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
func (m *migrator) Rollback() {
	err := m.db.Migrator().DropTable(
		&models.Template{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
