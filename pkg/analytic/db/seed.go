package db

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/analytic/service"
	user "github.com/axlle-com/blog/pkg/user/provider"
)

type seeder struct {
	analyticService *service.AnalyticService
	userProvider    user.UserProvider
}

func NewAnalyticSeeder(
	analyticService *service.AnalyticService,
	userProvider user.UserProvider,
) contracts.Seeder {
	return &seeder{
		analyticService: analyticService,
		userProvider:    userProvider,
	}
}

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	logger.Info("Database seeded Analytic successfully!")
}
