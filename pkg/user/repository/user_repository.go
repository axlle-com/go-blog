package repository

import (
	"github.com/axlle-com/blog/pkg/app/db"
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/user/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetAll() ([]*models.User, error)
	GetAllIds() ([]uint, error)
	GetByEmailWithRights(email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
	*app.Paginate
}

func NewUserRepo() UserRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) WithTx(tx *gorm.DB) UserRepository {
	newR := &repository{db: tx}
	return newR
}

func (r *repository) Create(user *models.User) error {
	user.Creating()
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
	user.Updating()
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.User{}).Pluck("id", &ids).Error; err != nil {
		return ids, err
	}
	return ids, nil
}
