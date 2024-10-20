package migrate

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func Migrate() {
	d := db.GetDB()

	err := d.AutoMigrate(
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
func Rollback() {
	d := db.GetDB()

	err := d.Migrator().DropTable(
		&models.Gallery{},
		&models.Image{},
		&models.GalleryHasResource{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
