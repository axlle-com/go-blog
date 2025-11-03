package dto

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

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

func (b InfoBlock) GetID() uint         { return b.ID }
func (b InfoBlock) GetUUID() uuid.UUID  { return parseUUID(b.UUID) }
func (b InfoBlock) GetTemplateID() uint { return b.TemplateID }

func (b InfoBlock) GetTemplateTitle() string { return b.Template }
func (b InfoBlock) GetTemplateName() string  { return b.Template }

func (b InfoBlock) GetTitle() string       { return b.Title }
func (b InfoBlock) GetDescription() string { return b.Description }
func (b InfoBlock) GetImage() string       { return b.Image }
func (b InfoBlock) GetMedia() string       { return b.Media }

func (b InfoBlock) GetPosition() string { return b.Position }

func (b InfoBlock) GetPositions() []string {
	return []string{
		"top",
		"bottom",
		"left",
		"right",
	}
}

func (b InfoBlock) GetSort() int        { return b.Sort }
func (b InfoBlock) GetRelationID() uint { return b.RelationID }

func (b InfoBlock) GetGalleries() []contract.Gallery {
	if len(b.Galleries) == 0 {
		return nil
	}
	out := make([]contract.Gallery, len(b.Galleries))
	for i := range b.Galleries {
		out[i] = b.Galleries[i] // dto.Gallery реализует contract.Gallery (см. ниже)
	}
	return out
}

// --- helpers ---

func parseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}
