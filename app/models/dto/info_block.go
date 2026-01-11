package dto

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type InfoBlock struct {
	ID          uint        `json:"id"`
	UUID        string      `json:"uuid"`
	TemplateID  *uint       `json:"template_id"`
	InfoBlockID *uint       `json:"info_block_id"`
	Template    string      `json:"template,omitempty"`
	Title       string      `json:"title"`
	Description *string     `json:"description,omitempty"`
	Image       *string     `json:"image,omitempty"`
	Media       *string     `json:"media,omitempty"`
	Position    string      `json:"position,omitempty"`
	Sort        int         `json:"sort"`
	RelationID  uint        `json:"relation_id,omitempty"`
	Galleries   []Gallery   `json:"galleries,omitempty"`
	InfoBlocks  []InfoBlock `json:"info_blocks,omitempty"`
}

func (i InfoBlock) GetID() uint          { return i.ID }
func (i InfoBlock) GetUUID() uuid.UUID   { return parseUUID(i.UUID) }
func (i InfoBlock) GetTemplateID() *uint { return i.TemplateID }

func (i InfoBlock) GetTemplateTitle() string { return i.Template }
func (i InfoBlock) GetTemplateName() string  { return i.Template }

func (i InfoBlock) GetTitle() string        { return i.Title }
func (i InfoBlock) GetDescription() *string { return i.Description }
func (i InfoBlock) GetImage() *string       { return i.Image }
func (i InfoBlock) GetMedia() *string       { return i.Media }

func (i InfoBlock) GetPosition() string { return i.Position }

func (i InfoBlock) GetPositions() []string {
	return []string{
		"top",
		"bottom",
		"left",
		"right",
	}
}

func (i InfoBlock) GetSort() int { return i.Sort }

func (i InfoBlock) GetRelationID() uint { return i.RelationID }

func (i InfoBlock) GetGalleries() []contract.Gallery {
	if len(i.Galleries) == 0 {
		return nil
	}

	out := make([]contract.Gallery, len(i.Galleries))
	for cnt := range i.Galleries {
		out[cnt] = i.Galleries[cnt]
	}
	return out
}

func (i InfoBlock) GetInfoBlockID() *uint {
	return i.InfoBlockID
}

func (i InfoBlock) GetInfoBlocks() []contract.InfoBlock {
	if len(i.InfoBlocks) == 0 {
		return nil
	}
	out := make([]contract.InfoBlock, 0, len(i.InfoBlocks))
	for _, ch := range i.InfoBlocks {
		out = append(out, ch)
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
