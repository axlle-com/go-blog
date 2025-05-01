package migrate

import (
	"gorm.io/gorm"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/models"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contracts.Migrator {
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

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_uuid ON users USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_has_users_user_uuid ON user_has_users USING hash (user_uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_has_users_relation_uuid ON user_has_users USING hash (relation_uuid);`)

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
