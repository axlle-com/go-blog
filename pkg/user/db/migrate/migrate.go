package migrate

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/user/models"
)

func Migrate() {
	d := db.GetDB()

	err := d.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
func Rollback() {
	d := db.GetDB()

	dropIntermediateTables()
	err := d.Migrator().DropTable(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}

func dropIntermediateTables() {
	d := db.GetDB()
	migrator := d.Migrator()
	intermediateTables := []string{
		"user_has_role",
		"user_has_permission",
		"role_has_permission",
	}
	for _, table := range intermediateTables {
		if err := migrator.DropTable(table); err != nil {
			fmt.Println("Error dropping table:", table, err)
			return
		}
		logger.Info(fmt.Sprintf("Dropped intermediate table:%s", table))
	}
}
