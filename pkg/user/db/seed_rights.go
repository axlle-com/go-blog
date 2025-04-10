package db

import (
	"github.com/axlle-com/blog/app/logger"
	. "github.com/axlle-com/blog/pkg/user/models"
)

func (s *seeder) seedPermissions() {
	permissions := []string{"create", "update", "delete", "read"}
	for _, name := range permissions {
		permission := Permission{
			Name: name,
		}
		err := s.permission.Create(&permission)
		if err != nil {
			logger.Errorf("Failed to create permission %v", err.Error())
		}
	}

	logger.Info("Database seeded Permission successfully!")
}

func (s *seeder) seedRoles() {
	roles := []string{"admin", "employee"}
	for _, name := range roles {
		role := Role{
			Name: name,
		}
		err := s.role.Create(&role)
		if err != nil {
			logger.Errorf("Failed to create role %v", err.Error())
		}
	}

	logger.Info("Database seeded Role successfully!")
}
