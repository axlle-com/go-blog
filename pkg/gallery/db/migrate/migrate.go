package migrate

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contracts.Migrator {
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

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_gallery_has_resources_resource_uuid ON gallery_has_resources USING hash (resource_uuid);`)

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
