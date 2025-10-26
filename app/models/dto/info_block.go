package dto

type InfoBlock struct {
	ID          uint      `json:"id"`
	UUID        string    `json:"uuid"`
	TemplateID  uint      `json:"template_id"`
	Template    string    `json:"template,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	Media       string    `json:"media,omitempty"`
	Position    string    `json:"position,omitempty"`
	Sort        int       `json:"sort"`
	RelationID  uint      `json:"relation_id,omitempty"`
	Galleries   []Gallery `json:"galleries,omitempty"`
}

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

type Image struct {
	ID          uint    `json:"id"`
	GalleryID   uint    `json:"gallery_id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Sort        int     `json:"sort"`
	File        string  `json:"file"`
}
