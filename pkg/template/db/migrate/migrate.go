package migrate

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/template/models"
)

func Migrate() {
	d := db.GetDB()

	err := d.AutoMigrate(
		&models.Template{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
func Rollback() {
	d := db.GetDB()

	err := d.Migrator().DropTable(
		&models.Template{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
