package contracts

type Seeder interface {
	Seed() error
	SeedTest(n int) error
}

type Migrator interface {
	Migrate() error
	Rollback() error
}
