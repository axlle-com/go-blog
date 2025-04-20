package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
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
		&models.PostTag{},
		&models.PostTagHasResource{},
	)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_uuid ON posts USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_alias ON posts USING hash (alias);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_url ON posts USING hash (url);`)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_categories_uuid ON post_categories USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_categories_alias ON post_categories USING hash (alias);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_categories_url ON post_categories USING hash (url);`)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_tags_alias ON post_tags USING hash (alias);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_tags_url ON post_tags USING hash (url);`)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_tag_has_resources_resource_uuid ON post_tag_has_resources USING hash (resource_uuid);`)

	if err != nil {
		logger.Fatal(err)
	}
}
func (m *migrator) Rollback() {
	err := m.db.Migrator().DropTable(
		&models.Post{},
		&models.PostCategory{},
		&models.PostTag{},
		&models.PostTagHasResource{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}
