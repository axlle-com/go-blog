package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contract.Migrator {
	return &migrator{db: db}
}

func (m *migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.CreateHashIndex("galleries", "uuid"))
	m.db.Exec(db.CreateHashIndex("gallery_has_resources", "resource_uuid"))

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)
	if err != nil {
		return err
	}

	return nil
}
