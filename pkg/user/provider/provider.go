package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	appProvider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/axlle-com/blog/pkg/user/service"
	"github.com/google/uuid"
)

func NewProvider(
	user repository.UserRepository,
	userService *service.UserService,
) appProvider.UserProvider {
	return &provider{
		userRepo:    user,
		userService: userService,
	}
}

type provider struct {
	userRepo    repository.UserRepository
	userService *service.UserService
}

func (p *provider) GetAll() []contract.User {
	all, err := p.userRepo.GetAll()
	if err == nil && len(all) > 0 {
		var users []contract.User
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

func (p *provider) GetByID(id uint) (contract.User, error) {
	t, err := p.userRepo.GetByID(id)
	if err == nil {
		return t, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetByUUID(uuid uuid.UUID) (contract.User, error) {
	t, err := p.userRepo.GetByUUID(uuid)
	if err == nil {
		return t, nil
	}

	logger.Error(err)
	return nil, err
}

func (p *provider) GetByIDs(ids []uint) ([]contract.User, error) {
	all, err := p.userRepo.GetByIDs(ids)
	if err == nil && len(all) > 0 {
		collection := make([]contract.User, 0, len(all))
		for _, item := range all {
			collection = append(collection, item)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contract.User, error) {
	all, err := p.userRepo.GetByIDs(ids)

	if err == nil && len(all) > 0 {
		users := make(map[uint]contract.User)
		for _, item := range all {
			users[item.ID] = item
		}

		return users, nil
	}

	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByUUIDs(uuids []uuid.UUID) (map[uuid.UUID]contract.User, error) {
	all, err := p.userRepo.GetByUUIDs(uuids)

	if err == nil && len(all) > 0 {
		users := make(map[uuid.UUID]contract.User)
		for _, item := range all {
			users[item.UUID] = item
		}

		return users, nil
	}

	logger.Error(err)
	return nil, err
}

func (p *provider) Create(user contract.User) (contract.User, error) {
	return p.userService.CreateFromInterface(user)
}
