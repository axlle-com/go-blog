package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/axlle-com/blog/app/models/contracts"
)

type Menu struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UUID        uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	TemplateID  *uint      `gorm:"index" json:"template_id"`
	Name        string     `gorm:"size:100" json:"name,omitempty"`
	IsPublished bool       `gorm:"default:true" json:"is_published,omitempty"`
	IsMain      bool       `gorm:"default:false" json:"IsMain,omitempty"`
	Ico         *string    `gorm:"size:255" json:"ico,omitempty"`
	Sort        int        `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt   *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Template contracts.Template `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
}

func (m *Menu) GetUUID() uuid.UUID {
	return m.UUID
}

func (m *Menu) GetTemplateName() string {
	if m.Template != nil {
		if m.Template.GetName() == "" {
			return fmt.Sprintf("%s.default", m.GetTable())
		}
		return fmt.Sprintf("%s.%s", m.GetTable(), m.Template.GetName())
	}
	return fmt.Sprintf("%s.default", m.GetTable())
}

func (m *Menu) GetName() string {
	return m.GetTable()
}

func (m *Menu) SetUUID() {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
}

func (m *Menu) GetTable() string {
	return "menus"
}

func (m *Menu) GetTemplateID() uint {
	var templateID uint
	if m.TemplateID != nil {
		templateID = *m.TemplateID
	}
	return templateID
}

func (m *Menu) Creating() {
	m.Saving()
}

func (m *Menu) Updating() {
	m.Saving()
}

func (m *Menu) Deleting() bool {
	return true
}

func (m *Menu) Saving() {
	m.SetUUID()
}
