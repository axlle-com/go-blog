package app

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	galleryAjax "github.com/axlle-com/blog/pkg/gallery/http/handlers/web"
	galleryProvider "github.com/axlle-com/blog/pkg/gallery/provider"
	galleryRepo "github.com/axlle-com/blog/pkg/gallery/repository"
	galleryService "github.com/axlle-com/blog/pkg/gallery/service"
	postAjax "github.com/axlle-com/blog/pkg/post/http/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/post/http/handlers/api"
	postWeb "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/axlle-com/blog/pkg/post/service"
	templateProvider "github.com/axlle-com/blog/pkg/template/provider"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
)

type Container struct {
	FileService  *file.Service
	FileProvider fileProvider.FileProvider

	GalleryResourceRepository galleryRepo.GalleryResourceRepository

	ImageRepo     galleryRepo.GalleryImageRepository
	ImageEvent    *galleryService.ImageEvent
	ImageService  *galleryService.ImageService
	ImageProvider galleryProvider.ImageProvider

	GalleryRepo     galleryRepo.GalleryRepository
	GalleryEvent    *galleryService.GalleryEvent
	GalleryService  *service.Service
	GalleryProvider galleryProvider.GalleryProvider

	PostRepo     repository.PostRepository
	PostService  *service.Service
	CategoryRepo repository.CategoryRepository

	TemplateProvider   templateProvider.TemplateProvider
	TemplateRepository templateRepository.TemplateRepository

	UserRepository userRepository.UserRepository
	UserProvider   userProvider.UserProvider

	AliasRepo     alias.AliasRepository
	AliasProvider alias.AliasProvider
}

func New() *Container {
	fService := file.NewService()
	fProvider := fileProvider.NewProvider(fService)

	iRepo := galleryRepo.NewImageRepo()
	iEvent := galleryService.NewImageEvent(fProvider)
	iService := galleryService.NewImageService(iRepo, iEvent)
	iProvider := galleryProvider.NewImageProvider(iRepo)

	rRepo := galleryRepo.NewResourceRepo()

	gRepo := galleryRepo.NewGalleryRepo()
	gEvent := galleryService.NewGalleryEvent(iService, rRepo)
	gService := galleryService.NewGalleryService(gRepo, gEvent, iService, rRepo)
	gProvider := galleryProvider.NewProvider(gRepo, gService)

	tRepo := templateRepository.NewTemplateRepo()
	tProvider := templateProvider.NewProvider(tRepo)

	uRepo := userRepository.NewUserRepo()
	uProvider := userProvider.NewProvider(uRepo)

	aRepo := alias.NewAliasRepo()
	aProvider := alias.NewProvider(aRepo)

	pRepo := repository.NewPostRepo()
	pService := service.NewService(pRepo, gProvider, fProvider, aProvider)
	cRepo := repository.NewCategoryRepo()

	return &Container{
		FileService:  fService,
		FileProvider: fProvider,

		GalleryResourceRepository: rRepo,

		ImageRepo:     iRepo,
		ImageEvent:    iEvent,
		ImageService:  iService,
		ImageProvider: iProvider,

		GalleryProvider: gProvider,
		GalleryRepo:     gRepo,
		GalleryService:  pService,
		GalleryEvent:    gEvent,

		PostRepo:     pRepo,
		PostService:  pService,
		CategoryRepo: cRepo,

		TemplateProvider:   tProvider,
		TemplateRepository: tRepo,

		UserRepository: uRepo,
		UserProvider:   uProvider,

		AliasRepo:     aRepo,
		AliasProvider: aProvider,
	}
}

func (c *Container) PostApiController() postApi.Controller {
	return postApi.New(
		c.PostService,
		c.PostRepo,
		c.CategoryRepo,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) PostController() postAjax.Controller {
	return postAjax.New(
		c.PostService,
		c.PostRepo,
		c.CategoryRepo,
		c.TemplateProvider,
		c.UserProvider,
	)
}

func (c *Container) PostWebController() postWeb.Controller {
	return postWeb.NewWebController(
		c.PostService,
		c.PostRepo,
		c.CategoryRepo,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) GalleryAjaxController() galleryAjax.Controller {
	return galleryAjax.New(
		c.GalleryRepo,
		c.ImageRepo,
		c.ImageService,
	)
}
