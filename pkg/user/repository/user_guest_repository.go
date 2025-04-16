package repository

import (
	"github.com/axlle-com/blog/app/db"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGuestRepository interface {
	WithTx(tx *gorm.DB) UserGuestRepository
	Create(user *models.UserGuest) error
	GetByID(id uint) (*models.UserGuest, error)
	GetByIDs(ids []uint) ([]*models.UserGuest, error)
	GetByUUIDs(uuids []uuid.UUID) ([]*models.UserGuest, error)
	GetByEmail(email string) (*models.UserGuest, error)
	Update(user *models.UserGuest) error
	Delete(id uint) error
	GetAll() ([]*models.UserGuest, error)
	GetAllIds() ([]uint, error)
	GetByEmailWithRights(email string) (*models.UserGuest, error)
}

type userGuestRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewUserGuestRepo() UserGuestRepository {
	return &userGuestRepository{db: db.GetDB()}
}

func (r *userGuestRepository) WithTx(tx *gorm.DB) UserGuestRepository {
	newR := &userGuestRepository{db: tx}
	return newR
}

func (r *userGuestRepository) Create(user *models.UserGuest) error {
	user.Creating()
	return r.db.Create(user).Error
}

func (r *userGuestRepository) GetByID(id uint) (*models.UserGuest, error) {
	var user models.UserGuest
	if err := r.db.Select(user.Fields()).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userGuestRepository) GetByIDs(ids []uint) ([]*models.UserGuest, error) {
	var users []*models.UserGuest
	if err := r.db.Select((&models.UserGuest{}).Fields()).Where("id IN (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userGuestRepository) GetByUUIDs(uuids []uuid.UUID) ([]*models.UserGuest, error) {
	var users []*models.UserGuest
	if err := r.db.Select((&models.UserGuest{}).Fields()).Where("uuid IN (?)", uuids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userGuestRepository) GetByEmail(email string) (*models.UserGuest, error) {
	var user models.UserGuest
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Roles.Permissions").Preload("Permissions").Find(&user)
	return &user, nil
}

func (r *userGuestRepository) GetByEmailWithRights(email string) (*models.UserGuest, error) {
	var user models.UserGuest
	if err := r.db.
		Preload("Roles").
		Preload("Permissions").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userGuestRepository) Update(user *models.UserGuest) error {
	user.Updating()
	return r.db.Save(user).Error
}

func (r *userGuestRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserGuest{}, id).Error
}

func (r *userGuestRepository) GetAll() ([]*models.UserGuest, error) {
	var users []*models.UserGuest
	if err := r.db.Select((&models.UserGuest{}).Fields()).Order("id ASC").Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func (r *userGuestRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.UserGuest{}).Pluck("id", &ids).Error; err != nil {
		return ids, err
	}
	return ids, nil
}
