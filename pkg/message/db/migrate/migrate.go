package migrate

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
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
		&models.Message{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_messages_uuid ON messages USING hash (user_uuid);`)

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.Message{},
	)

	if err != nil {
		return err
	}

	return nil
}
