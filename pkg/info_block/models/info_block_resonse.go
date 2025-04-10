package models

import (
	"github.com/axlle-com/blog/app/db"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
)

type InfoBlockResponse struct {
	ID          uint      `json:"id" form:"id" binding:"-"`
	UUID        uuid.UUID `json:"uuid" form:"uuid" binding:"-"`
	TemplateID  *uint     `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID      *uint     `json:"user_id" form:"user_id" binding:"omitempty"`
	Media       *string   `json:"media" form:"media" binding:"omitempty,max=255"`
	Title       string    `json:"title" form:"title" binding:"required,max=255"`
	Description *string   `json:"description" form:"description" binding:"omitempty"`
	Image       *string   `json:"image" form:"image" binding:"omitempty,max=255"`

	Template  contracts2.Template  `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
	User      contracts2.User      `gorm:"-" json:"user" form:"user" binding:"-" ignore:"true"`
	Galleries []contracts2.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`

	RelationID uint   `gorm:"relation_id" json:"relation_id" form:"relation_id" binding:"omitempty"`
	Sort       int    `gorm:"sort" json:"sort" form:"sort" binding:"omitempty"`
	Position   string `gorm:"position" json:"position" form:"position" binding:"omitempty"`
}

func (i *InfoBlockResponse) GetUUID() uuid.UUID {
	return i.UUID
}

func (i *InfoBlockResponse) GetPosition() string {
	return i.Position
}

func (i *InfoBlockResponse) GetPositions() []string {
	return []string{
		"top",
		"bottom",
		"left",
		"right",
	}
}

func (i *InfoBlockResponse) GetSort() int {
	return i.Sort
}

func (i *InfoBlockResponse) GetID() uint {
	return i.ID
}

func (i *InfoBlockResponse) GetTemplateID() uint {
	var templateID uint
	if i.TemplateID != nil {
		templateID = *i.TemplateID
	}
	return templateID
}

func (i *InfoBlockResponse) GetTemplateTitle() string {
	var title string
	if i.Template != nil {
		title = i.Template.GetTitle()
	}
	return title
}

func (i *InfoBlockResponse) UserLastName() string {
	var lastName string
	if i.User != nil {
		lastName = i.User.GetLastName()
	}
	return lastName
}

func (i *InfoBlockResponse) GetTitle() string {
	return i.Title
}

func (i *InfoBlockResponse) GetMedia() string {
	return *i.Media
}

func (i *InfoBlockResponse) GetDescription() string {
	return *i.Description
}

func (i *InfoBlockResponse) GetImage() string {
	return *i.Image
}

func (i *InfoBlockResponse) GetGalleries() []contracts2.Gallery {
	return i.Galleries
}

func (i *InfoBlockResponse) GetRelationID() uint {
	return i.RelationID
}

func (i *InfoBlockResponse) FromInterface(infoBlock contracts2.InfoBlock) {
	i.ID = infoBlock.GetID()
	i.UUID = infoBlock.GetUUID()
	i.TemplateID = db.UintPtr(infoBlock.GetTemplateID())
	i.Title = infoBlock.GetTitle()
	i.Description = db.StrPtr(infoBlock.GetDescription())
	i.Image = db.StrPtr(infoBlock.GetImage())
	i.Media = db.StrPtr(infoBlock.GetMedia())
}
