package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/axlle-com/blog/pkg/user/repository"
)

type UserProvider interface {
	GetAll() []*user.User
	GetAllIds() []uint
	GetByID(id uint) (*user.User, error)
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

func (p *provider) GetAll() []*user.User {
	all, err := p.userRepo.GetAll()
	if err == nil {
		return all
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

func (p *provider) GetByID(id uint) (*user.User, error) {
	t, err := p.userRepo.GetByID(id)
	if err == nil {
		return t, nil
	}
	logger.Error(err)
	return nil, err
}
