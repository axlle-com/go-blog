package api

import (
	apppPovider "github.com/axlle-com/blog/app/models/provider"
)

type Api struct {
	File      apppPovider.FileProvider
	Image     apppPovider.ImageProvider
	Gallery   apppPovider.GalleryProvider
	Blog      apppPovider.BlogProvider
	Template  apppPovider.TemplateProvider
	User      apppPovider.UserProvider
	Alias     apppPovider.AliasProvider
	InfoBlock apppPovider.InfoBlockProvider
	Analytic  apppPovider.AnalyticProvider
	Menu      apppPovider.MenuProvider
	Publisher apppPovider.PublisherProvider
}
