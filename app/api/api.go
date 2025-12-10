package api

import (
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/alias"
	analyticProvider "github.com/axlle-com/blog/pkg/analytic/provider"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	templateProvider "github.com/axlle-com/blog/pkg/template/provider"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
)

type Api struct {
	File      fileProvider.FileProvider
	Image     apppPovider.ImageProvider
	Gallery   apppPovider.GalleryProvider
	Blog      contract.BlogProvider
	Template  templateProvider.TemplateProvider
	User      userProvider.UserProvider
	Alias     alias.AliasProvider
	InfoBlock apppPovider.InfoBlockProvider
	Analytic  analyticProvider.AnalyticProvider
	Menu      apppPovider.MenuProvider
	Publisher apppPovider.PublisherProvider
}
