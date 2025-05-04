package db

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/user/models"
)

func (s *seeder) seedPermissions() error {
	permissions := []string{"create", "update", "delete", "read"}
	for _, name := range permissions {
		model, _ := s.permission.GetByName(name)
		if model != nil {
			continue
		}

		permission := models.Permission{
			Name: name,
		}
		err := s.permission.Create(&permission)
		if err != nil {
			return err
		}
	}

	logger.Info("[User][seeder][seedPermissions] Database seeded Permission successfully!")
	return nil
}

func (s *seeder) seedRoles() error {
	roles := []string{"superadmin", "admin", "employee"}
	for _, name := range roles {
		model, _ := s.role.GetByName(name)
		if model != nil {
			continue
		}

		role := models.Role{
			Name: name,
		}
		err := s.role.Create(&role)
		if err != nil {
			return err
		}
	}

	logger.Info("[User][seeder][seedRoles] Database seeded Role successfully!")
	return nil
}
