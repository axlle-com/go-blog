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
