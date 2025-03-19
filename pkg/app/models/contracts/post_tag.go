package contracts

type PostTag interface {
	GetID() uint
	GetTemplateID() uint
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetGalleries() []Gallery
}
