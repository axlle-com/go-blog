package contracts

type Seeder interface {
	Seed()
	SeedTest(n int)
}

type Migrator interface {
	Migrate()
	Rollback()
}
