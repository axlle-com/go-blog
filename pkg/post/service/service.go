package service

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/repository"
)

type Service struct {
	postRepo        repository.PostRepository
	galleryProvider gallery.GalleryProvider
	fileProvider    provider.FileProvider
	aliasProvider   alias.AliasProvider
}

func NewService(
	postRepo repository.PostRepository,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
) *Service {
	return &Service{
		postRepo:        postRepo,
		galleryProvider: galleryProvider,
		fileProvider:    fileProvider,
		aliasProvider:   aliasProvider,
	}
}
