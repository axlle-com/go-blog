package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetAll() ([]models.User, error)
	GetAllIds() ([]uint, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewRepository() UserRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) Create(user *models.User) error {
	user.SetPasswordHash()
	return r.db.Create(user).Error
}

func (r *repository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Roles.Permissions").Preload("Permissions").Find(&user)
	return &user, nil
}

func (r *repository) GetByEmailWithRights(email string) (*models.User, error) {
	var user models.User
	if err := r.db.
		Preload("Roles").
		Preload("Permissions").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user *models.User) error {
	user.SetPasswordHash()
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.User{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}
