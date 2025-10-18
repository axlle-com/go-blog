package model

import "github.com/google/uuid"

type Publisher struct {
	ID    uint   `json:"id" form:"id" binding:"omitempty"`
	UUID  string `json:"uuid" form:"uuid" binding:"omitempty"`
	URL   string `json:"url" form:"url" binding:"omitempty,max=1000"`
	Title string `json:"title" form:"title" binding:"omitempty,max=1000"`
}

func (p *Publisher) GetID() uint {
	return p.ID
}

func (p *Publisher) GetUUID() uuid.UUID {
	if parsed, err := uuid.Parse(p.UUID); err == nil {
		return parsed
	}

	return uuid.New()
}

func (p *Publisher) GetURL() string {
	return p.URL
}

func (p *Publisher) GetTitle() string {
	return p.Title
}

func (p *Publisher) GetTable() string {
	return ""
}
