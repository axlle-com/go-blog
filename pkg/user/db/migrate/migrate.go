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
	model := &models.User{}
	modelHas := &models.UserHasUser{}

	err := m.db.AutoMigrate(
		model,
		modelHas,
		&models.Role{},
		&models.Permission{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.HashIndex(model.GetTable(), "uuid"))
	m.db.Exec(db.HashIndex(modelHas.GetTable(), "user_uuid"))
	m.db.Exec(db.HashIndex(modelHas.GetTable(), "relation_uuid"))

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
	newMigrator := m.db.Migrator()
	intermediateTables := []string{
		"user_has_role",
		"user_has_permission",
		"role_has_permission",
	}

	for _, table := range intermediateTables {
		if err := newMigrator.DropTable(table); err != nil {
			return err
		}

		logger.Infof("[user][migrator][dropIntermediateTables] Dropped intermediate table:%s", table)
	}

	return nil
}
