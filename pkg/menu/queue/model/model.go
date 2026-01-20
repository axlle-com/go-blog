package model

import "github.com/google/uuid"

type Publisher struct {
	ID              uint    `json:"id" form:"id" binding:"omitempty"`
	UUID            string  `json:"uuid" form:"uuid" binding:"omitempty"`
	URL             string  `json:"url" form:"url" binding:"omitempty,max=1000"`
	Title           string  `json:"title" form:"title" binding:"omitempty,max=1000"`
	Image           *string `json:"image,omitempty"`
	MetaTitle       *string `json:"meta_title,omitempty"`
	MetaDescription *string `json:"meta_description,omitempty"`
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

func (p *Publisher) GetImage() string {
	if p.Image != nil {
		return *p.Image
	}
	return ""
}

func (p *Publisher) GetMetaTitle() string {
	if p.MetaTitle != nil {
		return *p.MetaTitle
	}
	return ""
}

func (p *Publisher) GetMetaDescription() string {
	if p.MetaDescription != nil {
		return *p.MetaDescription
	}
	return ""
}
