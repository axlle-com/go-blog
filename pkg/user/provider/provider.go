package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/repository"
)

type UserProvider interface {
	GetAll() []contracts.User
	GetAllIds() []uint
	GetByID(id uint) (contracts.User, error)
	GetByIDs(ids []uint) ([]contracts.User, error)
	GetMapByIDs(ids []uint) (map[uint]contracts.User, error)
}

func NewProvider(
	user repository.UserRepository,
) UserProvider {
	return &provider{
		userRepo: user,
	}
}

type provider struct {
	userRepo repository.UserRepository
}

func (p *provider) GetAll() []contracts.User {
	all, err := p.userRepo.GetAll()
	if err == nil && len(all) > 0 {
		var users []contracts.User
		for _, item := range all {
			users = append(users, item)
		}
		return users
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetAllIds() []uint {
	t, err := p.userRepo.GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetByID(id uint) (contracts.User, error) {
	t, err := p.userRepo.GetByID(id)
	if err == nil {
		return t, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetByIDs(ids []uint) ([]contracts.User, error) {
	all, err := p.userRepo.GetByIDs(ids)
	if err == nil && len(all) > 0 {
		collection := make([]contracts.User, 0, len(all))
		for _, item := range all {
			collection = append(collection, item)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contracts.User, error) {
	all, err := p.userRepo.GetByIDs(ids)

	if err == nil && len(all) > 0 {
		users := make(map[uint]contracts.User)
		for _, item := range all {
			users[item.ID] = item
		}

		return users, nil
	}

	logger.Error(err)
	return nil, err
}
