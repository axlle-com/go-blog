package models

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"time"
)

type Gallery struct {
	ID           uint            `gorm:"primary_key" json:"id"`
	Title        *string         `gorm:"size:255" json:"title"`
	Description  *string         `gorm:"type:text" json:"description"`
	Sort         int             `gorm:"default:0" json:"sort"`
	Image        *string         `gorm:"size:255;" json:"image"`
	URL          *string         `gorm:"size:255" json:"url"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
	DeletedAt    *time.Time      `gorm:"index" json:"deleted_at,omitempty"`
	GalleryImage []*GalleryImage `json:"images,omitempty"`
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

func (g *Gallery) GetImages() []contracts.GalleryImage {
	images := make([]contracts.GalleryImage, len(g.GalleryImage))
	for i, image := range g.GalleryImage {
		images[i] = image
	}
	return images
}

func (g *Gallery) Attach(r contracts.Resource) error {
	repo := NewGalleryResourceRepository()
	hasRepo, err := repo.GetByResourceAndID(r.GetID(), r.GetResource(), g.ID)
	if err != nil || hasRepo == nil {
		err = repo.Create(
			&GalleryHasResource{
				ResourceID: r.GetID(),
				Resource:   r.GetResource(),
				GalleryID:  g.ID,
			},
		)
	}
	return err
}

func (g *Gallery) Deleted() error {
	repo := NewGalleryResourceRepository()
	has, err := repo.GetByID(g.ID)
	if err == nil || has != nil {
		return repo.Delete(g.ID)
	}
	return err
}
