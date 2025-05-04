package contracts

import "gorm.io/gorm"

type Seeder interface {
	Seed() error
	SeedTest(n int) error
}

type Migrator interface {
	Migrate() error
	Rollback() error
}

type DB interface {
	GORM() *gorm.DB
	Close() error
}
