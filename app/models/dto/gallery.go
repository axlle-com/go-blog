package dto

import (
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type Gallery struct {
	ID           uint    `json:"id"`
	ResourceUUID string  `json:"resource_uuid"`
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	Sort         int     `json:"sort"`
	Position     string  `json:"position,omitempty"`
	Image        *string `json:"image,omitempty"`
	URL          *string `json:"url,omitempty"`
	Images       []Image `json:"images,omitempty"`
}

func (g Gallery) GetID() uint                { return g.ID }
func (g Gallery) GetResourceUUID() uuid.UUID { return parseUUID(g.ResourceUUID) }
func (g Gallery) GetTitle() *string          { return g.Title }
func (g Gallery) GetDescription() *string    { return g.Description }
func (g Gallery) GetSort() int               { return g.Sort }
func (g Gallery) GetPosition() string        { return g.Position }
func (g Gallery) GetImage() *string          { return g.Image }
func (g Gallery) GetURL() *string            { return g.URL }
func (g Gallery) GetDate() *time.Time        { return nil }
func (g Gallery) GetImages() []contract.Image {
	if len(g.Images) == 0 {
		return nil
	}
	out := make([]contract.Image, len(g.Images))
	for i := range g.Images {
		out[i] = g.Images[i]
	}
	return out
}

type Image struct {
	ID          uint    `json:"id"`
	GalleryID   uint    `json:"gallery_id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Sort        int     `json:"sort"`
	File        string  `json:"file"`
}

func (im Image) GetID() uint             { return im.ID }
func (im Image) GetGalleryID() uint      { return im.GalleryID }
func (im Image) GetTitle() *string       { return im.Title }
func (im Image) GetDescription() *string { return im.Description }
func (im Image) GetSort() int            { return im.Sort }
func (im Image) GetFile() string         { return im.File }
func (im Image) GetDate() *time.Time     { return nil }
func (im Image) GetGallery() contract.Gallery {
	return Gallery{ID: im.GalleryID}
}
