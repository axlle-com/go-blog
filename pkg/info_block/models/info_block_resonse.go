package models

import (
	"fmt"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
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

	Template  contract.Template  `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
	User      contract.User      `gorm:"-" json:"user" form:"user" binding:"-" ignore:"true"`
	Galleries []contract.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`

	Children     []*InfoBlockResponse `gorm:"-" json:"children" form:"children" binding:"-" ignore:"true"`
	RelationID   uint                 `gorm:"relation_id" json:"relation_id" form:"relation_id" binding:"omitempty"`
	ResourceUUID uuid.UUID            `gorm:"resource_uuid" json:"resource_uuid" form:"resource_uuid" binding:"omitempty"`
	Sort         int                  `gorm:"sort" json:"sort" form:"sort" binding:"omitempty"`
	Position     string               `gorm:"position" json:"position" form:"position" binding:"omitempty"`
}

func (i *InfoBlockResponse) GetUUID() uuid.UUID {
	return i.UUID
}

func (i *InfoBlockResponse) GetName() string {
	return (&InfoBlock{}).GetTable()
}

func (i *InfoBlockResponse) GetTemplateName() string {
	if i.Template != nil {
		return i.Template.GetFullName(i.GetName())
	}

	return fmt.Sprintf("%s.default", i.GetName())
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
	if i.Media != nil {
		return *i.Media
	}
	return ""
}

func (i *InfoBlockResponse) GetDescription() string {
	if i.Description != nil {
		return *i.Description
	}
	return ""
}

func (i *InfoBlockResponse) GetImage() string {
	return *i.Image
}

func (i *InfoBlockResponse) GetGalleries() []contract.Gallery {
	return i.Galleries
}

func (i *InfoBlockResponse) GetRelationID() uint {
	return i.RelationID
}

func (i *InfoBlockResponse) FromInterface(infoBlock contract.InfoBlock) {
	i.ID = infoBlock.GetID()
	i.UUID = infoBlock.GetUUID()
	i.TemplateID = db.UintPtr(infoBlock.GetTemplateID())
	i.Title = infoBlock.GetTitle()
	i.Description = db.StrPtr(infoBlock.GetDescription())
	i.Image = db.StrPtr(infoBlock.GetImage())
	i.Media = db.StrPtr(infoBlock.GetMedia())
}
