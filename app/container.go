package app

import (
	"context"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/cache"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/alias"
	provider2 "github.com/axlle-com/blog/pkg/analytic/provider"
	repository4 "github.com/axlle-com/blog/pkg/analytic/repository"
	service6 "github.com/axlle-com/blog/pkg/analytic/service"
	"github.com/axlle-com/blog/pkg/file"
	"github.com/axlle-com/blog/pkg/file/http"
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
	"github.com/axlle-com/blog/pkg/mailer"
	contracts2 "github.com/axlle-com/blog/pkg/message/contracts"
	web6 "github.com/axlle-com/blog/pkg/message/http/admin/handlers/ajax"
	web5 "github.com/axlle-com/blog/pkg/message/http/admin/handlers/web"
	web7 "github.com/axlle-com/blog/pkg/message/http/front/handlers/ajax"
	repository3 "github.com/axlle-com/blog/pkg/message/repository"
	service5 "github.com/axlle-com/blog/pkg/message/service"
	ajax2 "github.com/axlle-com/blog/pkg/post/http/admin/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/post/http/admin/handlers/api"
	web3 "github.com/axlle-com/blog/pkg/post/http/admin/handlers/web"
	web2 "github.com/axlle-com/blog/pkg/post/http/front/handlers/web"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/axlle-com/blog/pkg/post/service"
	"github.com/axlle-com/blog/pkg/queue"
	ajax3 "github.com/axlle-com/blog/pkg/template/http/admin/handlers/ajax"
	web4 "github.com/axlle-com/blog/pkg/template/http/admin/handlers/web"
	templateProvider "github.com/axlle-com/blog/pkg/template/provider"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	service4 "github.com/axlle-com/blog/pkg/template/service"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
	service3 "github.com/axlle-com/blog/pkg/user/service"
	"github.com/axlle-com/blog/pkg/view"
)

type Container struct {
	Queue contracts.Queue
	Cache contracts.Cache
	View  contracts.View

	FileService  *file.Service
	FileProvider fileProvider.FileProvider

	ImageRepo     galleryRepo.GalleryImageRepository
	ImageEvent    *galleryService.ImageEvent
	ImageService  *galleryService.ImageService
	ImageProvider galleryProvider.ImageProvider

	GalleryRepo         galleryRepo.GalleryRepository
	GalleryEvent        *galleryService.GalleryEvent
	GalleryService      *service.PostService
	GalleryProvider     galleryProvider.GalleryProvider
	GalleryResourceRepo galleryRepo.GalleryResourceRepository

	PostRepo          repository.PostRepository
	PostService       *service.PostService
	PostsService      *service.PostsService
	CategoryRepo      repository.CategoryRepository
	CategoriesService *service.CategoriesService
	CategoryService   *service.CategoryService

	TemplateProvider          templateProvider.TemplateProvider
	TemplateRepo              templateRepository.TemplateRepository
	TemplateService           *service4.TemplateService
	TemplateCollectionService *service4.TemplateCollectionService

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

	MessageRepo              contracts2.MessageRepository
	MessageService           *service5.MessageService
	MessageCollectionService *service5.MessageCollectionService
	MailService              *service5.MailService

	AnalyticRepo              repository4.AnalyticRepository
	AnalyticService           *service6.AnalyticService
	AnalyticCollectionService *service6.AnalyticCollectionService
	AnalyticProvider          provider2.AnalyticProvider
}

func New(ctx context.Context) *Container {
	newQueue := queue.NewQueue()
	newQueue.StartWorkers(ctx, 4)

	newCache := cache.NewCache()
	newView := view.NewView(config.Config())

	mailerInterface := mailer.NewMailer(newQueue)

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

	uRepo := userRepository.NewUserRepo()
	uService := service3.NewUserService(uRepo)
	uaService := service3.NewAuthService(uService)
	uProvider := userProvider.NewProvider(uRepo, uService)

	tRepo := templateRepository.NewTemplateRepo()
	tProvider := templateProvider.NewProvider(tRepo)
	tService := service4.NewTemplateService(tRepo, uProvider)
	tCollectionService := service4.NewTemplateCollectionService(tService, tRepo, uProvider)

	mRepo := repository3.NewMessageRepo()
	mService := service5.NewMessageService(mRepo, uProvider)
	mcService := service5.NewMessageCollectionService(mRepo, mService, uProvider)
	mailService := service5.NewMailService(mService, mcService, uProvider, mailerInterface, newQueue)

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

	analyticRepo := repository4.NewAnalyticRepo()
	analyticService := service6.NewAnalyticService(analyticRepo, uProvider)
	analyticCollectionService := service6.NewAnalyticCollectionService(analyticRepo, analyticService, uProvider)
	analyticProvider := provider2.NewAnalyticProvider(analyticService, analyticCollectionService)

	return &Container{
		Queue: newQueue,
		Cache: newCache,
		View:  newView,

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

		TemplateProvider:          tProvider,
		TemplateRepo:              tRepo,
		TemplateService:           tService,
		TemplateCollectionService: tCollectionService,

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

		MessageRepo:              mRepo,
		MessageService:           mService,
		MessageCollectionService: mcService,
		MailService:              mailService,

		AnalyticRepo:              analyticRepo,
		AnalyticService:           analyticService,
		AnalyticCollectionService: analyticCollectionService,
		AnalyticProvider:          analyticProvider,
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
		c.View,
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) TemplateWebController() web4.TemplateWebController {
	return web4.NewTemplateWebController(
		c.TemplateService,
		c.TemplateCollectionService,
		c.UserProvider,
	)
}

func (c *Container) TemplateController() ajax3.TemplateController {
	return ajax3.NewTemplateController(
		c.TemplateService,
		c.TemplateCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageController() web5.MessageWebController {
	return web5.NewMessageWebController(
		c.MessageService,
		c.MessageCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageAjaxController() web6.MessageController {
	return web6.NewMessageController(
		c.MessageService,
		c.MessageCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageFrontController() web7.MessageController {
	return web7.NewMessageController(
		c.MailService,
	)
}

func (c *Container) UserFrontController() user.Controller {
	return user.NewUserWebController(
		c.UserService,
		c.UserAuthService,
		c.Cache,
	)
}

func (c *Container) FileController() http.Controller {
	return http.NewFileController(
		c.FileService,
	)
}
