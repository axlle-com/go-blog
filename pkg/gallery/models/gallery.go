package models

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"time"
)

type Gallery struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Title       *string    `gorm:"size:255" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Sort        int        `gorm:"default:0" json:"sort"`
	Image       *string    `gorm:"size:255;" json:"image"`
	URL         *string    `gorm:"size:255" json:"url"`
	CreatedAt   *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Images      []*Image   `json:"images,omitempty"`
}

func (g *Gallery) GetID() uint {
	return g.ID
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
