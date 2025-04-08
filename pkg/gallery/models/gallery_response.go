package models

import (
	"github.com/google/uuid"
	"time"

	"github.com/axlle-com/blog/pkg/app/models/contracts"
)

type GalleryResponse struct {
	ID          uint       `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Image       *string    `json:"image"`
	URL         *string    `json:"url"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	Sort         int       `json:"sort" form:"sort" binding:"omitempty"`
	Position     string    `json:"position" form:"position" binding:"omitempty"`
	ResourceUUID uuid.UUID `json:"uuid" form:"uuid" binding:"-"`
	Images       []*Image  `json:"images,omitempty" gorm:"foreignKey:GalleryID;references:ID"`
}

func (g *GalleryResponse) GetID() uint {
	return g.ID
}

func (g *GalleryResponse) GetResourceUUID() uuid.UUID {
	return g.ResourceUUID
}

func (g *GalleryResponse) GetTitle() *string {
	return g.Title
}

func (g *GalleryResponse) GetDescription() *string {
	return g.Description
}

func (g *GalleryResponse) GetSort() int {
	return g.Sort
}

func (g *GalleryResponse) GetPosition() string {
	return g.Position
}

func (g *GalleryResponse) GetImage() *string {
	return g.Image
}

func (g *GalleryResponse) GetURL() *string {
	return g.URL
}

func (g *GalleryResponse) GetDate() *time.Time {
	return g.CreatedAt
}

func (g *GalleryResponse) GetImages() []contracts.Image {
	images := make([]contracts.Image, len(g.Images))
	for i, image := range g.Images {
		images[i] = image
	}
	return images
}
