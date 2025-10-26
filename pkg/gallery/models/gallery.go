package models

import (
	"time"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
)

type Gallery struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	UUID        uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	Title       *string    `gorm:"size:255" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Image       *string    `gorm:"size:255;" json:"image"`
	URL         *string    `gorm:"size:255" json:"url"`
	CreatedAt   *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// @todo delete?
	Sort         int       `gorm:"-" json:"sort" form:"sort" binding:"omitempty"`
	Position     string    `gorm:"-" json:"position" form:"position" binding:"omitempty"`
	ResourceUUID uuid.UUID `gorm:"-" json:"resource_uuid" form:"resource_uuid" binding:"-"`
	Images       []*Image  `json:"images,omitempty"`
}

func (g *Gallery) GetID() uint {
	return g.ID
}

func (g *Gallery) GetResourceUUID() uuid.UUID {
	return g.ResourceUUID
}

func (g *Gallery) GetTitle() *string {
	return g.Title
}

func (g *Gallery) GetDescription() *string {
	return g.Description
}

func (g *Gallery) GetSort() int {
	return g.Sort
}

func (g *Gallery) GetPosition() string {
	return g.Position
}

func (g *Gallery) GetImage() *string {
	return g.Image
}

func (g *Gallery) GetURL() *string {
	return g.URL
}

func (g *Gallery) GetDate() *time.Time {
	return g.CreatedAt
}

func (g *Gallery) GetImages() []contracts.Image {
	images := make([]contracts.Image, len(g.Images))
	for i, image := range g.Images {
		images[i] = image
	}
	return images
}

func (g *Gallery) Saving() {
	g.SetUUID()
}

func (g *Gallery) SetUUID() {
	if g.UUID == uuid.Nil {
		g.UUID = uuid.New()
	}
}
