package contracts

type InfoBlock interface {
	GetID() uint
	GetTemplateID() uint
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetGalleries() []Gallery
}
