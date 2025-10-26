package dto

import (
	"github.com/axlle-com/blog/app/models/contracts"
)

func MapInfoBlock(src contracts.InfoBlock) InfoBlock {
	if src == nil {
		return InfoBlock{}
	}
	return InfoBlock{
		ID:          src.GetID(),
		UUID:        src.GetUUID().String(),
		TemplateID:  src.GetTemplateID(),
		Template:    src.GetTemplateName(),
		Title:       src.GetTitle(),
		Description: src.GetDescription(),
		Image:       src.GetImage(),
		Media:       src.GetMedia(),
		Position:    src.GetPosition(),
		Sort:        src.GetSort(),
		RelationID:  src.GetRelationID(),
		Galleries:   MapGalleries(src.GetGalleries()),
	}
}

func MapInfoBlocks(list []contracts.InfoBlock) []InfoBlock {
	if len(list) == 0 {
		return nil
	}
	out := make([]InfoBlock, 0, len(list))
	for _, ib := range list {
		out = append(out, MapInfoBlock(ib))
	}

	return out
}

func MapGallery(src contracts.Gallery) Gallery {
	if src == nil {
		return Gallery{}
	}
	return Gallery{
		ID:           src.GetID(),
		ResourceUUID: src.GetResourceUUID().String(),
		Title:        copyStrPtr(src.GetTitle()),
		Description:  copyStrPtr(src.GetDescription()),
		Sort:         src.GetSort(),
		Position:     src.GetPosition(),
		Image:        copyStrPtr(src.GetImage()),
		URL:          copyStrPtr(src.GetURL()),
		Images:       MapImages(src.GetImages()),
	}
}

func MapGalleries(list []contracts.Gallery) []Gallery {
	if len(list) == 0 {
		return nil
	}
	out := make([]Gallery, 0, len(list))
	for _, g := range list {
		out = append(out, MapGallery(g))
	}
	return out
}

func MapImage(src contracts.Image) Image {
	if src == nil {
		return Image{}
	}
	return Image{
		ID:          src.GetID(),
		GalleryID:   src.GetGalleryID(),
		Title:       copyStrPtr(src.GetTitle()),
		Description: copyStrPtr(src.GetDescription()),
		Sort:        src.GetSort(),
		File:        src.GetFile(),
	}
}

func MapImages(list []contracts.Image) []Image {
	if len(list) == 0 {
		return nil
	}
	out := make([]Image, 0, len(list))
	for _, im := range list {
		out = append(out, MapImage(im))
	}
	return out
}

func copyStrPtr(p *string) *string {
	if p == nil {
		return nil
	}
	s := *p
	return &s
}
