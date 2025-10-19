package models

type GalleryRequest struct {
	ID          uint            `json:"id" form:"id" binding:"omitempty"`
	Title       *string         `json:"title" form:"title" binding:"omitempty"`
	Description *string         `json:"description" form:"description" binding:"omitempty"`
	Sort        int             `json:"sort" form:"sort" binding:"omitempty"`
	Position    string          `json:"position" form:"position" binding:"omitempty"`
	Image       *string         `json:"image" form:"image" binding:"omitempty"`
	URL         *string         `json:"url" form:"url" binding:"omitempty"`
	Images      []*ImageRequest `json:"images" form:"images" binding:"omitempty"`
}

type ImageRequest struct {
	ID           uint    `json:"id" form:"id" binding:"omitempty"`
	GalleryID    uint    `json:"gallery_id" form:"gallery_id" binding:"omitempty"`
	OriginalName string  `json:"original_name" form:"original_name" binding:"omitempty"`
	File         string  `json:"file" form:"file" binding:"omitempty"`
	Title        *string `json:"title" form:"title" binding:"omitempty"`
	Description  *string `json:"description" form:"description" binding:"omitempty"`
	Sort         int     `json:"sort" form:"sort" binding:"omitempty"`
}

type InfoBlockRequest struct {
	ID          uint    `json:"id" form:"id" binding:"omitempty"`
	Title       string  `json:"title" form:"title" binding:"omitempty"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
	UserID      *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Media       *string `json:"media" form:"media" binding:"omitempty,max=255"`
	Sort        int     `json:"sort" form:"sort" binding:"omitempty"`
	Position    string  `json:"position" form:"position" binding:"omitempty"`

	GalleryID  *uint `json:"gallery_id" form:"gallery_id" binding:"omitempty"`
	TemplateID *uint `json:"template_id" form:"template_id" binding:"omitempty"`
	RelationID uint  `json:"relation_id" form:"relation_id" binding:"omitempty"`
}
