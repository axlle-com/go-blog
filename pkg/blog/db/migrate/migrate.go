package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
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
		&models.Post{},
		&models.PostCategory{},
		&models.PostTag{},
		&models.PostTagHasResource{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.CreateHashIndex("posts", "uuid"))
	m.db.Exec(db.CreateHashIndex("posts", "alias"))
	m.db.Exec(db.CreateHashIndex("posts", "url"))

	m.db.Exec(db.CreateHashIndex("post_categories", "uuid"))
	m.db.Exec(db.CreateHashIndex("post_categories", "alias"))
	m.db.Exec(db.CreateHashIndex("post_categories", "url"))

	m.db.Exec(db.CreateHashIndex("post_tags", "name"))
	m.db.Exec(db.CreateHashIndex("post_tags", "alias"))
	m.db.Exec(db.CreateHashIndex("post_tags", "url"))

	m.db.Exec(db.CreateHashIndex("post_tag_has_resources", "resource_uuid"))

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.Post{},
		&models.PostCategory{},
		&models.PostTag{},
		&models.PostTagHasResource{},
	)

	if err != nil {
		return err
	}

	return nil
}
