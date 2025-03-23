package app

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	galleryAjax "github.com/axlle-com/blog/pkg/gallery/http/handlers/web"
	galleryProvider "github.com/axlle-com/blog/pkg/gallery/provider"
	galleryRepo "github.com/axlle-com/blog/pkg/gallery/repository"
	galleryService "github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/info_block/provider"
	repository2 "github.com/axlle-com/blog/pkg/info_block/repository"
	service2 "github.com/axlle-com/blog/pkg/info_block/service"
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

	GalleryResourceRepo galleryRepo.GalleryResourceRepository

	ImageRepo     galleryRepo.GalleryImageRepository
	ImageEvent    *galleryService.ImageEvent
	ImageService  *galleryService.ImageService
	ImageProvider galleryProvider.ImageProvider

	GalleryRepo     galleryRepo.GalleryRepository
	GalleryEvent    *galleryService.GalleryEvent
	GalleryService  *service.Service
	GalleryProvider galleryProvider.GalleryProvider

	PostRepo          repository.PostRepository
	PostService       *service.Service
	CategoryRepo      repository.CategoryRepository
	CategoriesService *service.CategoriesService

	TemplateProvider templateProvider.TemplateProvider
	TemplateRepo     templateRepository.TemplateRepository

	UserRepo     userRepository.UserRepository
	UserProvider userProvider.UserProvider

	AliasRepo     alias.AliasRepository
	AliasProvider alias.AliasProvider

	InfoBlockHasResourceRepo repository2.InfoBlockHasResourceRepository
	InfoBlockRepo            repository2.InfoBlockRepository
	InfoBlockService         *service2.InfoBlockService
	InfoBlockProvider        provider.InfoBlockProvider

	PostTagRepo         repository.PostTagRepository
	PostTagResourceRepo repository.PostTagResourceRepository
	PostTagService      *service.PostTagService
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
	cService := service.NewCategoryService(cRepo, tProvider, uProvider)

	ibhrRepo := repository2.NewResourceRepo()
	ibRepo := repository2.NewInfoBlockRepo()
	ibService := service2.NewInfoBlockService(ibRepo, ibhrRepo)
	ibProvider := provider.NewProvider(ibRepo, ibService)

	ptRepo := repository.NewPostTagRepo()
	ptrRepo := repository.NewResourceRepo()
	ptService := service.NewPostTagService(ptRepo, ptrRepo)

	return &Container{
		FileService:  fService,
		FileProvider: fProvider,

		GalleryResourceRepo: rRepo,

		ImageRepo:     iRepo,
		ImageEvent:    iEvent,
		ImageService:  iService,
		ImageProvider: iProvider,

		GalleryProvider: gProvider,
		GalleryRepo:     gRepo,
		GalleryService:  pService,
		GalleryEvent:    gEvent,

		PostRepo:          pRepo,
		PostService:       pService,
		CategoryRepo:      cRepo,
		CategoriesService: cService,

		TemplateProvider: tProvider,
		TemplateRepo:     tRepo,

		UserRepo:     uRepo,
		UserProvider: uProvider,

		AliasRepo:     aRepo,
		AliasProvider: aProvider,

		InfoBlockHasResourceRepo: ibhrRepo,
		InfoBlockRepo:            ibRepo,
		InfoBlockService:         ibService,
		InfoBlockProvider:        ibProvider,

		PostTagRepo:         ptRepo,
		PostTagResourceRepo: ptrRepo,
		PostTagService:      ptService,
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

func (c *Container) PostCategoryWebController() postWeb.ControllerCategory {
	return postWeb.NewWebControllerCategory(
		c.CategoryRepo,
		c.CategoriesService,
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
