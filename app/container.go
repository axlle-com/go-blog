package app

import (
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/file"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	galleryAjax "github.com/axlle-com/blog/pkg/gallery/http/handlers/web"
	galleryProvider "github.com/axlle-com/blog/pkg/gallery/provider"
	galleryRepo "github.com/axlle-com/blog/pkg/gallery/repository"
	galleryService "github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/info_block/http/admin/handlers/ajax"
	"github.com/axlle-com/blog/pkg/info_block/http/admin/handlers/web"
	"github.com/axlle-com/blog/pkg/info_block/provider"
	repository2 "github.com/axlle-com/blog/pkg/info_block/repository"
	service2 "github.com/axlle-com/blog/pkg/info_block/service"
	ajax2 "github.com/axlle-com/blog/pkg/post/http/admin/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/post/http/admin/handlers/api"
	web3 "github.com/axlle-com/blog/pkg/post/http/admin/handlers/web"
	web2 "github.com/axlle-com/blog/pkg/post/http/front/handlers/web"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/axlle-com/blog/pkg/post/service"
	templateProvider "github.com/axlle-com/blog/pkg/template/provider"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
	service3 "github.com/axlle-com/blog/pkg/user/service"
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
	GalleryService  *service.PostService
	GalleryProvider galleryProvider.GalleryProvider

	PostRepo          repository.PostRepository
	PostService       *service.PostService
	PostsService      *service.PostsService
	CategoryRepo      repository.CategoryRepository
	CategoriesService *service.CategoriesService
	CategoryService   *service.CategoryService

	TemplateProvider templateProvider.TemplateProvider
	TemplateRepo     templateRepository.TemplateRepository

	UserRepo        userRepository.UserRepository
	UserProvider    userProvider.UserProvider
	UserService     *service3.UserService
	UserAuthService *service3.AuthService

	AliasRepo     alias.AliasRepository
	AliasProvider alias.AliasProvider

	InfoBlockHasResourceRepo   repository2.InfoBlockHasResourceRepository
	InfoBlockRepo              repository2.InfoBlockRepository
	InfoBlockService           *service2.InfoBlockService
	InfoBlockCollectionService *service2.InfoBlockCollectionService
	InfoBlockProvider          provider.InfoBlockProvider

	PostTagRepo         repository.PostTagRepository
	PostTagResourceRepo repository.PostTagResourceRepository
	PostTagService      *service.PostTagService
}

func New() *Container {
	fileService := file.NewService()
	fileProv := fileProvider.NewProvider(fileService)

	iRepo := galleryRepo.NewImageRepo()
	iEvent := galleryService.NewImageEvent(fileProv)
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
	uService := service3.NewUserService(uRepo)
	uaService := service3.NewAuthService(uService)

	aRepo := alias.NewAliasRepo()
	aProvider := alias.NewProvider(aRepo)

	pRepo := repository.NewPostRepo()

	cRepo := repository.NewCategoryRepo()

	ibhrRepo := repository2.NewResourceRepo()
	ibRepo := repository2.NewInfoBlockRepo()

	ibcService := service2.NewInfoBlockCollectionService(ibRepo, ibhrRepo, gProvider, tProvider, uProvider)
	ibService := service2.NewInfoBlockService(ibRepo, ibcService, ibhrRepo, gProvider, tProvider, uProvider)
	ibProvider := provider.NewProvider(ibService, ibcService)

	ptRepo := repository.NewPostTagRepo()
	ptrRepo := repository.NewResourceRepo()
	ptService := service.NewPostTagService(ptRepo, ptrRepo)

	csService := service.NewCategoriesService(cRepo, aProvider, gProvider, tProvider, uProvider)
	cService := service.NewCategoryService(cRepo, aProvider, gProvider, fileProv, ibProvider)

	pService := service.NewPostService(pRepo, csService, cService, gProvider, fileProv, aProvider, ibProvider)
	psService := service.NewPostsService(pRepo, csService, cService, gProvider, fileProv, aProvider, uProvider, tProvider, ibProvider)

	return &Container{
		FileService:  fileService,
		FileProvider: fileProv,

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
		PostsService:      psService,
		CategoryRepo:      cRepo,
		CategoriesService: csService,
		CategoryService:   cService,

		TemplateProvider: tProvider,
		TemplateRepo:     tRepo,

		UserRepo:        uRepo,
		UserProvider:    uProvider,
		UserService:     uService,
		UserAuthService: uaService,

		AliasRepo:     aRepo,
		AliasProvider: aProvider,

		InfoBlockHasResourceRepo:   ibhrRepo,
		InfoBlockRepo:              ibRepo,
		InfoBlockService:           ibService,
		InfoBlockCollectionService: ibcService,
		InfoBlockProvider:          ibProvider,

		PostTagRepo:         ptRepo,
		PostTagResourceRepo: ptrRepo,
		PostTagService:      ptService,
	}
}

func (c *Container) PostApiController() postApi.Controller {
	return postApi.New(
		c.PostService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) PostController() ajax2.Controller {
	return ajax2.New(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) PostWebController() web3.Controller {
	return web3.NewWebController(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) CategoryWebController() web3.ControllerCategory {
	return web3.NewWebControllerCategory(
		c.CategoriesService,
		c.CategoryService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) CategoryController() ajax2.CategoryController {
	return ajax2.NewCategoryController(
		c.CategoriesService,
		c.CategoryService,
		c.TemplateProvider,
		c.UserProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) GalleryAjaxController() galleryAjax.Controller {
	return galleryAjax.New(
		c.GalleryRepo,
		c.ImageRepo,
		c.ImageService,
	)
}

func (c *Container) InfoBlockController() ajax.InfoBlockController {
	return ajax.NewInfoBlockController(
		c.InfoBlockService,
		c.InfoBlockCollectionService,
		c.TemplateProvider,
		c.UserProvider,
	)
}

func (c *Container) InfoBlockWebController() web.InfoBlockWebController {
	return web.NewInfoBlockWebController(
		c.InfoBlockService,
		c.InfoBlockCollectionService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) PostFrontWebController() web2.PostController {
	return web2.NewFrontWebController(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}
