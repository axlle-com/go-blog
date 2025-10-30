package migrate

import (
	"github.com/axlle-com/blog/app/models/contract"
)

type seeder struct {
	seeders []contract.Seeder
}

func NewSeeder(arg ...contract.Seeder) contract.Seeder {
	m := &seeder{}
	m.seeders = append(m.seeders, arg...)
	return m
}

func (s *seeder) Seed() error {
	for _, item := range s.seeders {
		if err := item.Seed(); err != nil {
			return err
		}
	}
	return nil
}

func (s *seeder) SeedTest(n int) error {
	for _, item := range s.seeders {
		if err := item.SeedTest(n); err != nil {
			return err
		}
	}
	return nil
}
