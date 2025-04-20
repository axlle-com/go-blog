package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
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
		&models.Message{},
	)

	if err != nil {
		logger.Fatal(err)
	}
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_messages_uuid ON messages USING hash (user_uuid);`)
}
func (m *migrator) Rollback() {
	err := m.db.Migrator().DropTable(
		&models.Message{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
