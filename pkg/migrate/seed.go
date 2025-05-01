package migrate

import (
	"github.com/axlle-com/blog/app/models/contracts"
)

type seeder struct {
	seeders []contracts.Seeder
}

func NewSeeder(arg ...contracts.Seeder) contracts.Seeder {
	m := &seeder{}
	m.seeders = append(m.seeders, arg...)
	return m
}

func (s *seeder) Seed() error {
	for _, mig := range s.seeders {
		if err := mig.Seed(); err != nil {
			return err
		}
	}
	return nil
}

func (s *seeder) SeedTest(n int) error {
	for _, mig := range s.seeders {
		if err := mig.SeedTest(n); err != nil {
			return err
		}
	}
	return nil
}
