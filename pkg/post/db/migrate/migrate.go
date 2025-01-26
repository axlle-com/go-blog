package migrate

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/post/models"
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
		&models.Post{},
		&models.PostCategory{},
	)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_uuid ON posts USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_categories_uuid ON post_categories USING hash (uuid);`)

	if err != nil {
		logger.Fatal(err)
	}
}
func (m *migrator) Rollback() {
	err := m.db.Migrator().DropTable(
		&models.Post{},
		&models.PostCategory{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
