package db

import (
	. "github.com/axlle-com/blog/pkg/user/models"
	. "github.com/axlle-com/blog/pkg/user/repository"
	"log"
)

func SeedPermissions() {
	permissions := []string{"create", "update"}
	for _, name := range permissions {
		permission := Permission{
			Name: name,
		}
		err := NewPermissionRepository().Create(&permission)
		if err != nil {
			log.Printf("Failed to create permission %v", err.Error())
		}
	}

	log.Println("Database seeded Permission successfully!")
}

func SeedRoles() {
	roles := []string{"admin", "employee"}
	for _, name := range roles {
		role := Role{
			Name: name,
		}
		err := NewRoleRepository().Create(&role)
		if err != nil {
			log.Printf("Failed to create role %v", err.Error())
		}
	}

	log.Println("Database seeded Role successfully!")
}
