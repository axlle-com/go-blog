package dto

import (
	"github.com/axlle-com/blog/app/models/contract"
)

func MapInfoBlock(src contract.InfoBlock) InfoBlock {
	if src == nil {
		return InfoBlock{}
	}

	out := InfoBlock{
		ID:          src.GetID(),
		UUID:        src.GetUUID().String(),
		TemplateID:  src.GetTemplateID(),
		InfoBlockID: src.GetInfoBlockID(),
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

	children := src.GetInfoBlocks()
	if len(children) > 0 {
		out.InfoBlocks = make([]InfoBlock, 0, len(children))
		for _, ch := range children {
			out.InfoBlocks = append(out.InfoBlocks, MapInfoBlock(ch))
		}
	}

	return out
}

func MapInfoBlocks(list []contract.InfoBlock) []InfoBlock {
	if len(list) == 0 {
		return nil
	}

	out := make([]InfoBlock, 0, len(list))
	for _, ib := range list {
		out = append(out, MapInfoBlock(ib))
	}

	return out
}

func MapGallery(src contract.Gallery) Gallery {
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

func MapGalleries(list []contract.Gallery) []Gallery {
	if len(list) == 0 {
		return nil
	}

	out := make([]Gallery, 0, len(list))
	for _, g := range list {
		out = append(out, MapGallery(g))
	}

	return out
}

func MapImage(src contract.Image) Image {
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

func MapImages(list []contract.Image) []Image {
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
