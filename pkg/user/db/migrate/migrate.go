package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/models"
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
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
func (m *migrator) Rollback() {
	m.dropIntermediateTables()
	err := m.db.Migrator().DropTable(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}

func (m *migrator) dropIntermediateTables() {
	migrator := m.db.Migrator()
	intermediateTables := []string{
		"user_has_role",
		"user_has_permission",
		"role_has_permission",
	}
	for _, table := range intermediateTables {
		if err := migrator.DropTable(table); err != nil {
			logger.Errorf("Error dropping table: %v, %v", table, err)
			return
		}
		logger.Infof("Dropped intermediate table:%s", table)
	}
}
