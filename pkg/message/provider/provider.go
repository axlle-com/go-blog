package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/service"
)

type MessageProvider interface {
	GetAll() []contracts.Message
	GetByID(id uint) (contracts.Message, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contracts.Message, error)
	GetMapByIDs(ids []uint) (map[uint]contracts.Message, error)
}

func NewMessageProvider(
	messageService service.MessageService,
	messageCollectionService service.MessageCollectionService,
) MessageProvider {
	return &provider{
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
	}
}

type provider struct {
	messageService           service.MessageService
	messageCollectionService service.MessageCollectionService
}

func (p *provider) GetAll() []contracts.Message {
	all, err := p.messageCollectionService.GetAll()
	if err == nil {
		collection := make([]contracts.Message, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetAllIds() []uint {
	t, err := p.messageCollectionService.GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetByID(id uint) (contracts.Message, error) {
	model, err := p.messageService.GetByID(id)
	if err == nil {
		return model, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetByIDs(ids []uint) ([]contracts.Message, error) {
	all, err := p.messageCollectionService.GetByIDs(ids)
	if err == nil {
		collection := make([]contracts.Message, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contracts.Message, error) {
	all, err := p.messageCollectionService.GetByIDs(ids)
	if err == nil {
		collection := make(map[uint]contracts.Message, len(all))
		for _, template := range all {
			collection[template.ID] = template
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}
