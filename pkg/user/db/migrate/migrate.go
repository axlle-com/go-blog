package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/user/models"
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
		&models.User{},
		&models.UserHasUser{},
		&models.Role{},
		&models.Permission{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.HashIndex("users", "uuid"))
	m.db.Exec(db.HashIndex("user_has_users", "user_uuid"))
	m.db.Exec(db.HashIndex("user_has_users", "relation_uuid"))

	return nil
}
func (m *migrator) Rollback() error {
	err := m.dropIntermediateTables()
	if err != nil {
		return err
	}

	return m.db.Migrator().DropTable(
		&models.User{},
		&models.UserHasUser{},
		&models.Role{},
		&models.Permission{},
	)
}

func (m *migrator) dropIntermediateTables() error {
	migrator := m.db.Migrator()
	intermediateTables := []string{
		"user_has_role",
		"user_has_permission",
		"role_has_permission",
	}
	for _, table := range intermediateTables {
		if err := migrator.DropTable(table); err != nil {
			logger.Errorf("Error dropping table: %v, %v", table, err)
			return err
		}
		logger.Infof("Dropped intermediate table:%s", table)
	}
	return nil
}
