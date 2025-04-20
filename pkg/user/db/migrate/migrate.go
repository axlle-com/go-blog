package migrate

import (
	"gorm.io/gorm"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/models"
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
		&models.UserHasUser{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(err)
	}

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_uuid ON users USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_has_users_user_uuid ON user_has_users USING hash (user_uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_has_users_relation_uuid ON user_has_users USING hash (relation_uuid);`)
}
func (m *migrator) Rollback() {
	m.dropIntermediateTables()
	err := m.db.Migrator().DropTable(
		&models.User{},
		&models.UserHasUser{},
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
