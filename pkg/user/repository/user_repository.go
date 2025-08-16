package repository

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Create(user *models.User) error
	Attach(userHasUser *models.UserHasUser) error
	GetByID(id uint) (*models.User, error)
	GetByUUID(uuid uuid.UUID) (*models.User, error)
	GetByIDs(ids []uint) ([]*models.User, error)
	GetByUUIDs(uuids []uuid.UUID) ([]*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetRelation(userUUID, relationUUID uuid.UUID) (*models.UserHasUser, error)
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

func NewUserRepo(db *gorm.DB) UserRepository {
	return &repository{db: db}
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
	if err := r.db.Select(user.Fields()).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByUUID(uuid uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Select(user.Fields()).Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByIDs(ids []uint) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Select((&models.User{}).Fields()).Where("id IN (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetByUUIDs(uuids []uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Select((&models.User{}).Fields()).Where("uuid IN (?)", uuids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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

func (r *repository) Attach(userHasUser *models.UserHasUser) error {
	return r.db.Create(userHasUser).Error
}

func (r *repository) GetRelation(userUUID, relationUUID uuid.UUID) (*models.UserHasUser, error) {
	var userHasUser models.UserHasUser
	if err := r.db.
		Where("user_uuid = ?", userUUID).
		Where("relation_uuid = ?", relationUUID).
		First(&userHasUser).Error; err != nil {
		return nil, err
	}
	return &userHasUser, nil
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Select((&models.User{}).Fields()).Order("id ASC").Find(&users).Error; err != nil {
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
