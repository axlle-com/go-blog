package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/user/repository"
)

type User interface {
	GetAll() []*models.User
	GetAllIds() []uint
}

func Provider() User {
	return &provider{}
}

type provider struct {
}

func (p *provider) GetAll() []*models.User {
	all, err := repository.NewRepo().GetAll()
	if err == nil {
		return all
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetAllIds() []uint {
	t, err := repository.NewRepo().GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}
