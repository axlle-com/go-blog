package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
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
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_gallery_has_resources_resource_uuid ON gallery_has_resources USING hash (resource_uuid);`)

	if err != nil {
		logger.Fatal(err)
	}
}
func (m *migrator) Rollback() {
	err := m.db.Migrator().DropTable(
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
