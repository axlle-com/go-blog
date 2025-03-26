package models

type GalleryRequest struct {
	ID          string          `json:"id" form:"id" binding:"omitempty"`
	Title       string          `json:"title" form:"title" binding:"omitempty"`
	Description string          `json:"description" form:"description" binding:"omitempty"`
	Sort        string          `json:"sort" form:"sort" binding:"omitempty"`
	Image       string          `json:"image" form:"image" binding:"omitempty"`
	URL         string          `json:"url" form:"url" binding:"omitempty"`
	Images      []*ImageRequest `json:"images" form:"images" binding:"omitempty"`
}

type ImageRequest struct {
	ID           string `json:"id" form:"id" binding:"omitempty"`
	GalleryID    string `json:"gallery_id" form:"gallery_id" binding:"omitempty"`
	OriginalName string `json:"original_name" form:"original_name" binding:"omitempty"`
	File         string `json:"file" form:"file" binding:"omitempty"`
	Title        string `json:"title" form:"title" binding:"omitempty"`
	Description  string `json:"description" form:"description" binding:"omitempty"`
	Sort         string `json:"sort" form:"sort" binding:"omitempty"`
}
