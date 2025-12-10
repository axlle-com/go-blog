package models

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

func FromContract(publisher contract.Publisher) contract.Publisher {
	return &Publisher{
		ID:    publisher.GetID(),
		UUID:  publisher.GetUUID(),
		URL:   publisher.GetURL(),
		Title: publisher.GetTitle(),
		Table: publisher.GetTable(),
	}
}

type Publisher struct {
	ID    uint      `json:"id"`
	UUID  uuid.UUID `json:"uuid" form:"uuid" binding:"-"`
	URL   string    `json:"url"`
	Title string    `json:"title"`
	Table string    `json:"table"`
}

func (p *Publisher) GetTable() string {
	return p.Table
}

func (p *Publisher) GetID() uint {
	return p.ID
}

func (p *Publisher) GetUUID() uuid.UUID {
	return p.UUID
}

func (p *Publisher) GetURL() string {
	return p.URL
}

func (p *Publisher) GetTitle() string {
	return p.Title
}
