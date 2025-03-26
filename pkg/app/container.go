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
	uService := service3.NewUserService(uRepo)
	uaService := service3.NewAuthService(uService)

	aRepo := alias.NewAliasRepo()
	aProvider := alias.NewProvider(aRepo)

	pRepo := repository.NewPostRepo()

	cRepo := repository.NewCategoryRepo()
	csService := service.NewCategoriesService(cRepo, aProvider, gProvider, tProvider, uProvider)
	cService := service.NewCategoryService(cRepo, aProvider, gProvider, fProvider)

	pService := service.NewPostService(pRepo, csService, cService, gProvider, fProvider, aProvider)
	psService := service.NewPostsService(pRepo, csService, cService, gProvider, fProvider, aProvider)

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
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) PostController() postAjax.Controller {
	return postAjax.New(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
	)
}

func (c *Container) PostWebController() postWeb.Controller {
	return postWeb.NewWebController(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) CategoryWebController() postWeb.ControllerCategory {
	return postWeb.NewWebControllerCategory(
		c.CategoriesService,
		c.CategoryService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) CategoryController() postAjax.CategoryController {
	return postAjax.NewCategoryController(
		c.CategoriesService,
		c.CategoryService,
		c.TemplateProvider,
		c.UserProvider,
	)
}

func (c *Container) GalleryAjaxController() galleryAjax.Controller {
	return galleryAjax.New(
		c.GalleryRepo,
		c.ImageRepo,
		c.ImageService,
	)
}
