package models

import (
	"github.com/google/uuid"
	"time"

	"github.com/axlle-com/blog/pkg/app/models/contracts"
)

type InfoBlock struct {
	ID          uint       `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	UUID        uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	TemplateID  *uint      `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	UserID      *uint      `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	Media       *string    `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title       string     `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	Description *string    `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	Image       *string    `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	CreatedAt   *time.Time `gorm:"index" json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" form:"updated_at" binding:"-" ignore:"true"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-" ignore:"true"`

	Template  contracts.Template  `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
	User      contracts.User      `gorm:"-" json:"user" form:"user" binding:"-" ignore:"true"`
	Galleries []contracts.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`

	Sort        int                   `gorm:"-" json:"sort" form:"sort" binding:"omitempty"`
	Position    string                `gorm:"-" json:"position" form:"position" binding:"omitempty"`
	HasResource *InfoBlockHasResource `gorm:"-" json:"has_resource" form:"has_resource" binding:"omitempty"`
}

func (i *InfoBlock) GetUUID() uuid.UUID {
	return i.UUID
}

func (i *InfoBlock) GetPosition() string {
	return i.Position
}

func (i *InfoBlock) GetPositions() []string {
	return []string{
		"top",
		"bottom",
		"left",
		"right",
	}
}

func (i *InfoBlock) GetSort() int {
	return i.Sort
}

func (i *InfoBlock) SetUUID() {
	if i.UUID == uuid.Nil {
		i.UUID = uuid.New()
	}
}

func (i *InfoBlock) GetID() uint {
	return i.ID
}

func (i *InfoBlock) GetTemplateID() uint {
	var templateID uint
	if i.TemplateID != nil {
		templateID = *i.TemplateID
	}
	return templateID
}

func (i *InfoBlock) GetRelationID() uint {
	var id uint
	if i.HasResource != nil {
		id = i.HasResource.ID
	}
	return id
}

func (i *InfoBlock) GetTemplateTitle() string {
	var title string
	if i.Template != nil {
		title = i.Template.GetTitle()
	}
	return title
}

func (i *InfoBlock) UserLastName() string {
	var lastName string
	if i.User != nil {
		lastName = i.User.GetLastName()
	}
	return lastName
}

func (i *InfoBlock) AdminURL() string {
	return "/admin/info-blocks"
}

func (i *InfoBlock) Date() string {
	if i.CreatedAt == nil {
		return ""
	}
	return i.CreatedAt.Format("02.01.2006 15:04:05")
}

func (i *InfoBlock) GetTable() string {
	return "info_blocks"
}

func (i *InfoBlock) GetTitle() string {
	return i.Title
}

func (i *InfoBlock) GetMedia() string {
	return *i.Media
}

func (i *InfoBlock) GetDescription() string {
	return *i.Description
}

func (i *InfoBlock) GetImage() string {
	return *i.Image
}

func (i *InfoBlock) GetGalleries() []contracts.Gallery {
	return i.Galleries
}

func (i *InfoBlock) Creating() {
	i.Saving()
}

func (i *InfoBlock) Updating() {
	i.Saving()
}

func (i *InfoBlock) Deleting() bool {
	return true
}

func (i *InfoBlock) Saving() {
	i.SetUUID()
}
