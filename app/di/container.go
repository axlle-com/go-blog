package di

import (
	"time"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/app/service/cache"
	"github.com/axlle-com/blog/app/service/disk"
	i18nsvc "github.com/axlle-com/blog/app/service/i18n"
	"github.com/axlle-com/blog/app/service/mailer"
	mailerQueue "github.com/axlle-com/blog/app/service/mailer/queue"
	"github.com/axlle-com/blog/app/service/migrate"
	"github.com/axlle-com/blog/app/service/minify"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/app/service/scheduler"
	"github.com/axlle-com/blog/app/service/storage"
	"github.com/axlle-com/blog/app/service/view"
	"github.com/axlle-com/blog/pkg/alias"
	analyticMigrate "github.com/axlle-com/blog/pkg/analytic/db/migrate"
	analyticProvider "github.com/axlle-com/blog/pkg/analytic/provider"
	analyticQueue "github.com/axlle-com/blog/pkg/analytic/queue"
	analyticRepo "github.com/axlle-com/blog/pkg/analytic/repository"
	analyticService "github.com/axlle-com/blog/pkg/analytic/service"
	postDB "github.com/axlle-com/blog/pkg/blog/db"
	postMigrate "github.com/axlle-com/blog/pkg/blog/db/migrate"
	postAjax "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/api"
	postAdminWeb "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/web"
	postFrontWeb "github.com/axlle-com/blog/pkg/blog/http/front/handlers/web"
	"github.com/axlle-com/blog/pkg/blog/provider"
	postQueue "github.com/axlle-com/blog/pkg/blog/queue"
	postRepo "github.com/axlle-com/blog/pkg/blog/repository"
	postService "github.com/axlle-com/blog/pkg/blog/service"
	fileMigrate "github.com/axlle-com/blog/pkg/file/db/migrate"
	fileAdminWeb "github.com/axlle-com/blog/pkg/file/http"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	fileQueue "github.com/axlle-com/blog/pkg/file/queue"
	fileRepo "github.com/axlle-com/blog/pkg/file/repository"
	fileService "github.com/axlle-com/blog/pkg/file/service"
	galleryMigrate "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	galleryAjax "github.com/axlle-com/blog/pkg/gallery/http/handlers/web"
	galleryProvider "github.com/axlle-com/blog/pkg/gallery/provider"
	galleryRepo "github.com/axlle-com/blog/pkg/gallery/repository"
	galleryService "github.com/axlle-com/blog/pkg/gallery/service"
	infoBlockDB "github.com/axlle-com/blog/pkg/info_block/db"
	infoBlockMigrate "github.com/axlle-com/blog/pkg/info_block/db/migrate"
	infoBlockAdminAjax "github.com/axlle-com/blog/pkg/info_block/http/admin/handlers/ajax"
	infoBlockAdminWeb "github.com/axlle-com/blog/pkg/info_block/http/admin/handlers/web"
	infoBlockProvider "github.com/axlle-com/blog/pkg/info_block/provider"
	infoBlockQueue "github.com/axlle-com/blog/pkg/info_block/queue"
	infoBlockRepo "github.com/axlle-com/blog/pkg/info_block/repository"
	infoBlockService "github.com/axlle-com/blog/pkg/info_block/service"
	menuDB "github.com/axlle-com/blog/pkg/menu/db"
	menuMigrate "github.com/axlle-com/blog/pkg/menu/db/migrate"
	menuAdminAjax "github.com/axlle-com/blog/pkg/menu/http/handlers/ajax"
	menuAdminWeb "github.com/axlle-com/blog/pkg/menu/http/handlers/web"
	menuProvider "github.com/axlle-com/blog/pkg/menu/provider"
	menuQueue "github.com/axlle-com/blog/pkg/menu/queue"
	menuRepository "github.com/axlle-com/blog/pkg/menu/repository"
	menuService "github.com/axlle-com/blog/pkg/menu/service"
	messageContracts "github.com/axlle-com/blog/pkg/message/contracts"
	messageDB "github.com/axlle-com/blog/pkg/message/db"
	messageMigrate "github.com/axlle-com/blog/pkg/message/db/migrate"
	messageAdminAjax "github.com/axlle-com/blog/pkg/message/http/admin/handlers/ajax"
	messageAdminWeb "github.com/axlle-com/blog/pkg/message/http/admin/handlers/web"
	messageFrontWeb "github.com/axlle-com/blog/pkg/message/http/front/handlers/ajax"
	messageQueue "github.com/axlle-com/blog/pkg/message/queue"
	messageRepo "github.com/axlle-com/blog/pkg/message/repository"
	messageService "github.com/axlle-com/blog/pkg/message/service"
	publisherAjax "github.com/axlle-com/blog/pkg/publisher/http/admin/handlers/ajax"
	publisherProvider "github.com/axlle-com/blog/pkg/publisher/provider"
	publisherService "github.com/axlle-com/blog/pkg/publisher/service"
	settingsDB "github.com/axlle-com/blog/pkg/settings/db"
	settingsMigrate "github.com/axlle-com/blog/pkg/settings/db/migrate"
	settingsProvider "github.com/axlle-com/blog/pkg/settings/provider"
	settingsRepo "github.com/axlle-com/blog/pkg/settings/repository"
	settingsService "github.com/axlle-com/blog/pkg/settings/service"
	templateDB "github.com/axlle-com/blog/pkg/template/db"
	templateMigrate "github.com/axlle-com/blog/pkg/template/db/migrate"
	templateAdminAjax "github.com/axlle-com/blog/pkg/template/http/admin/handlers/ajax"
	templateAdminWeb "github.com/axlle-com/blog/pkg/template/http/admin/handlers/web"
	templateProvider "github.com/axlle-com/blog/pkg/template/provider"
	templateRepo "github.com/axlle-com/blog/pkg/template/repository"
	templateService "github.com/axlle-com/blog/pkg/template/service"
	userDB "github.com/axlle-com/blog/pkg/user/db"
	userMigrate "github.com/axlle-com/blog/pkg/user/db/migrate"
	userFrontWeb "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	usersQueue "github.com/axlle-com/blog/pkg/user/queue"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
	usersService "github.com/axlle-com/blog/pkg/user/service"
	"github.com/gin-contrib/sessions/redis"
)

type Container struct {
	Config      contract.Config
	Queue       contract.Queue
	Cache       contract.Cache
	View        contract.View
	Scheduler   contract.Scheduler
	Disk        contract.DiskService
	Minifier    contract.Minifier
	SeedService contract.SeedService
	Store       redis.Store
	I18n        *i18nsvc.Service

	FileUploadService     *fileService.UploadService
	FileService           *fileService.Service
	FileCollectionService *fileService.CollectionService
	FileProvider          apppPovider.FileProvider

	ImageRepo     galleryRepo.GalleryImageRepository
	ImageEvent    *galleryService.ImageEvent
	ImageService  *galleryService.ImageService
	ImageProvider apppPovider.ImageProvider

	GalleryRepo         galleryRepo.GalleryRepository
	GalleryEvent        *galleryService.GalleryEvent
	GalleryService      *galleryService.GalleryService
	GalleryProvider     apppPovider.GalleryProvider
	GalleryResourceRepo galleryRepo.GalleryResourceRepository

	PostRepo          postRepo.PostRepository
	PostService       *postService.PostService
	PostsService      *postService.PostCollectionService
	PostProvider      apppPovider.BlogProvider
	CategoryRepo      postRepo.CategoryRepository
	CategoriesService *postService.CategoryCollectionService
	CategoryService   *postService.CategoryService

	TemplateProvider          apppPovider.TemplateProvider
	TemplateRepo              templateRepo.TemplateRepository
	TemplateService           *templateService.Service
	TemplateCollectionService *templateService.CollectionService

	UserRepo        userRepository.UserRepository
	UserProvider    apppPovider.UserProvider
	UserService     *usersService.UserService
	UserAuthService *usersService.AuthService

	UserRoleRepo       userRepository.RoleRepository
	UserPermissionRepo userRepository.PermissionRepository

	AliasRepo     alias.AliasRepository
	AliasProvider apppPovider.AliasProvider

	InfoBlockHasResourceRepo   infoBlockRepo.InfoBlockHasResourceRepository
	InfoBlockRepo              infoBlockRepo.InfoBlockRepository
	InfoBlockService           *infoBlockService.Service
	InfoBlockCollectionService *infoBlockService.CollectionService
	InfoBlockProvider          apppPovider.InfoBlockProvider

	PostTagRepo              postRepo.PostTagRepository
	PostTagResourceRepo      postRepo.PostTagResourceRepository
	PostTagService           *postService.TagService
	PostTagCollectionService *postService.TagCollectionService

	MessageRepo              messageContracts.MessageRepository
	MessageService           *messageService.MessageService
	MessageCollectionService *messageService.MessageCollectionService
	MailService              *messageService.MailService

	AnalyticRepo              analyticRepo.AnalyticRepository
	AnalyticService           *analyticService.Service
	AnalyticCollectionService *analyticService.CollectionService
	AnalyticProvider          apppPovider.AnalyticProvider

	MenuRepo                  menuRepository.MenuRepository
	MenuService               *menuService.MenuService
	MenuCollectionService     *menuService.MenuCollectionService
	MenuItemRepo              menuRepository.MenuItemRepository
	MenuItemService           *menuService.MenuItemService
	MenuItemCollectionService *menuService.MenuItemCollectionService
	MenuProvider              apppPovider.MenuProvider

	SettingsRepo    settingsRepo.Repository
	SettingsService *settingsService.Service

	Migrator contract.Migrator
	Seeder   contract.Seeder

	// Controllers
	AdminWebPostController       postAdminWeb.PostController
	AdminAjaxPostController      postAjax.PostController
	AdminApiPostController       postApi.Controller
	AdminWebCategoryController   postAdminWeb.CategoryController
	AdminAjaxCategoryController  postAjax.CategoryController
	AdminWebTagController        postAdminWeb.TagController
	AdminAjaxTagController       postAjax.TagController
	AdminWebInfoBlockController  infoBlockAdminWeb.InfoBlockWebController
	AdminAjaxInfoBlockController infoBlockAdminAjax.InfoBlockController
	AdminWebTemplateController   templateAdminWeb.TemplateWebController
	AdminAjaxTemplateController  templateAdminAjax.TemplateController
	AdminWebMessageController    messageAdminWeb.MessageWebController
	AdminAjaxMessageController   messageAdminAjax.MessageController
	AdminAjaxGalleryController   galleryAjax.Controller
	AdminWebMenuController       menuAdminWeb.Controller
	AdminAjaxMenuController      menuAdminAjax.ControllerMenu
	AdminAjaxMenuItemController  menuAdminAjax.ControllerMenuItem
	AdminWebFileController       fileAdminWeb.Controller
	FrontWebPostController       postFrontWeb.BlogController
	FrontAjaxMessageController   messageFrontWeb.MessageController
	FrontWebUserController       userFrontWeb.Controller
	AdminAjaxPublisherController publisherAjax.PublisherController

	Api *api.Api
}

func NewContainer(cfg contract.Config, db contract.DB) *Container {
	newQueue := queue.NewQueue()
	newCache := cache.NewCache()
	newDisk := disk.NewDiskService(cfg)
	newMinifier := minify.NewWebMinifier(cfg, newDisk)
	newStore := models.Store(cfg)
	newSeedService := migrate.NewSeedService(db, newDisk)

	newMailer := mailer.NewMailer(cfg, newQueue)
	newView := view.NewView(cfg, newDisk, newMinifier)

	newFileRepo := fileRepo.NewFileRepo(db.PostgreSQL())
	newFileService := fileService.NewService(newFileRepo)
	newStorageService := storage.NewLocalStorageService(cfg)
	uploadService := fileService.NewUploadService(newFileService, newStorageService)
	fileCollectionService := fileService.NewCollectionService(newFileRepo, newFileService, uploadService)
	fileProv := fileProvider.NewFileProvider(uploadService, newFileService, fileCollectionService)

	newUserRepo := userRepository.NewUserRepo(db.PostgreSQL())
	newRoleRepo := userRepository.NewRoleRepo(db.PostgreSQL())
	newPermissionRepo := userRepository.NewPermissionRepo(db.PostgreSQL())
	newUserService := usersService.NewUserService(newUserRepo, newRoleRepo, newPermissionRepo)
	newAuthService := usersService.NewAuthService(newUserService)
	newUserProvider := userProvider.NewProvider(newUserRepo, newUserService)

	newTemplateRepo := templateRepo.NewTemplateRepo(db.PostgreSQL())
	newTemplateProvider := templateProvider.NewProvider(newTemplateRepo)

	newAliasRepo := alias.NewAliasRepo(db.PostgreSQL())
	newAliasProvider := alias.NewAliasProvider(newAliasRepo)

	newImageRepo := galleryRepo.NewImageRepo(db.PostgreSQL())
	newImageProvider := galleryProvider.NewImageProvider(newImageRepo)

	newResourceRepo := galleryRepo.NewResourceRepo(db.PostgreSQL())
	newGalleryRepo := galleryRepo.NewGalleryRepo(db.PostgreSQL())

	newApi := &api.Api{
		File:        fileProv,
		Image:       newImageProvider,
		Gallery:     nil, // Will be set later
		Blog:        nil, // Will be set later
		Template:    newTemplateProvider,
		User:        newUserProvider,
		Alias:       newAliasProvider,
		InfoBlock:   nil, // Will be set later
		Analytic:    nil, // Will be set later
		Menu:        nil, // Will be set later
		Publisher:   nil, // Will be set later
		CompanyInfo: nil, // Will be set later
	}

	newImageEvent := galleryService.NewImageEvent(newApi)
	newImageService := galleryService.NewImageService(newImageRepo, newImageEvent, newApi)
	newGalleryEvent := galleryService.NewGalleryEvent(newQueue, newImageService, newGalleryRepo, newResourceRepo)
	newImageEvent.SetGalleryEvent(newGalleryEvent)
	newGalleryService := galleryService.NewGalleryService(newGalleryRepo, newGalleryEvent, newImageService, newResourceRepo, newApi)
	newGalleryProvider := galleryProvider.NewProvider(newGalleryRepo, newGalleryService)
	newApi.Gallery = newGalleryProvider

	newTemplateService := templateService.NewService(newApi, newTemplateRepo)
	newTemplateCollectionService := templateService.NewCollectionService(newApi, newTemplateRepo, newTemplateService)

	newMessageRepo := messageRepo.NewMessageRepo(db.PostgreSQL())
	newMessageService := messageService.NewMessageService(newMessageRepo, newApi)
	newMessageCollectionService := messageService.NewMessageCollectionService(newMessageRepo, newMessageService, newApi)
	newMailService := messageService.NewMailService(cfg, newQueue, newMessageService, newMessageCollectionService, newApi)

	newSettingsRepo := settingsRepo.NewRepository(db.PostgreSQL())
	newSettingsService := settingsService.NewService(newCache, newSettingsRepo, 60*24*time.Minute)
	newCompanyInfoService := settingsService.NewCompanyInfoService(newCache, newSettingsRepo, newSettingsService, 60*24*time.Minute)

	newApi.CompanyInfo = settingsProvider.NewProvider(newCompanyInfoService)

	newPostRepo := postRepo.NewPostRepo(db.PostgreSQL())
	newCategoryRepo := postRepo.NewCategoryRepo(db.PostgreSQL())

	newInfoBlockHasResourceRepo := infoBlockRepo.NewResourceRepo(db.PostgreSQL())
	newInfoBlockRepo := infoBlockRepo.NewInfoBlockRepo(db.PostgreSQL())

	newBlockCollectionAggregateService := infoBlockService.NewCollectionAggregateService(newApi, newInfoBlockRepo, newInfoBlockHasResourceRepo)
	newBlockCollectionService := infoBlockService.NewCollectionService(newApi, newInfoBlockRepo, newInfoBlockHasResourceRepo, newBlockCollectionAggregateService)
	newBlockEventService := infoBlockService.NewEventService(newQueue, newInfoBlockRepo)
	newBlockAggregateService := infoBlockService.NewAggregateService(newApi)
	newBlockService := infoBlockService.NewService(newApi, newInfoBlockRepo, newBlockCollectionService, newInfoBlockHasResourceRepo, newBlockEventService, newBlockAggregateService)
	newBlockProvider := infoBlockProvider.NewProvider(newBlockService, newBlockCollectionService)

	newApi.InfoBlock = newBlockProvider

	postTagRepo := postRepo.NewPostTagRepo(db.PostgreSQL())
	postTagResourceRepo := postRepo.NewResourceRepo(db.PostgreSQL())
	postTagService := postService.NewTagService(postTagRepo, postTagResourceRepo, newApi)
	postTagCollectionService := postService.NewTagCollectionService(postTagService, postTagRepo, postTagResourceRepo, newApi)

	newCategoriesService := postService.NewCategoryCollectionService(newCategoryRepo, newApi)
	categoryService := postService.NewCategoryService(newApi, newCategoryRepo)

	newPostAggregateService := postService.NewPostAggregateService(newApi, newPostRepo, newCategoriesService, categoryService, postTagCollectionService)
	newPostService := postService.NewPostService(newQueue, newApi, newPostRepo, newPostAggregateService, newCategoriesService, categoryService, postTagCollectionService)
	newPostCollectionService := postService.NewPostCollectionService(newPostRepo, newCategoriesService, categoryService, newApi)
	newCategoryAggregateService := postService.NewCategoryAggregateService(newApi, newCategoryRepo, newPostCollectionService)
	categoryService.SetAggregateService(newCategoryAggregateService)
	newBlogProvider := provider.NewBlogProvider(newPostService, newPostCollectionService, newCategoriesService, postTagCollectionService)

	newApi.Blog = newBlogProvider

	newAnalyticRepo := analyticRepo.NewAnalyticRepo(db.PostgreSQL())
	newAnalyticService := analyticService.NewService(newApi, newAnalyticRepo)
	analyticCollectionService := analyticService.NewCollectionService(newApi, newAnalyticRepo, newAnalyticService)
	newAnalyticProvider := analyticProvider.NewAnalyticProvider(newAnalyticService, analyticCollectionService)

	newApi.Analytic = newAnalyticProvider

	menuRepo := menuRepository.NewMenuRepo(db.PostgreSQL())
	menuItemRepo := menuRepository.NewMenuItemRepo(db.PostgreSQL())
	menuItemService := menuService.NewMenuItemService(menuItemRepo)
	menuItemAggregateService := menuService.NewMenuItemAggregateService(menuItemRepo, newApi)
	menuItemCollectionService := menuService.NewMenuItemCollectionService(menuItemRepo, menuItemService)
	newMenuService := menuService.NewMenuService(menuRepo, menuItemService, menuItemCollectionService, menuItemAggregateService)
	newMenuCollectionService := menuService.NewMenuCollectionService(menuRepo, newMenuService, newApi)
	newMenuProvider := menuProvider.NewMenuProvider(newView, newMenuService)

	newApi.Menu = newMenuProvider

	newPublisherService := publisherService.NewCollectionService(newApi)

	newApi.Publisher = publisherProvider.NewPublisherProvider(newApi)

	newI18n := i18nsvc.New(cfg, newDisk)

	menuSeeder := menuDB.NewMenuSeeder(cfg, newDisk, newSeedService, newApi, menuRepo, menuItemRepo)
	settingsSeeder := settingsDB.NewSettingsSeeder(cfg, newDisk, newSeedService, newSettingsService)

	userMigrator := userMigrate.NewMigrator(db.PostgreSQL())
	postMigrator := postMigrate.NewMigrator(db.PostgreSQL())
	infoBlockMigrator := infoBlockMigrate.NewMigrator(db.PostgreSQL())
	galleryMigrator := galleryMigrate.NewMigrator(db.PostgreSQL())
	templateMigrator := templateMigrate.NewMigrator(db.PostgreSQL())
	analyticMigrator := analyticMigrate.NewMigrator(db.PostgreSQL())
	messageMigrator := messageMigrate.NewMigrator(db.PostgreSQL())
	fileMigrator := fileMigrate.NewMigrator(db.PostgreSQL())
	menuMigrator := menuMigrate.NewMigrator(db.PostgreSQL())
	settingsMigrator := settingsMigrate.NewMigrator(db.PostgreSQL())

	userSeeder := userDB.NewSeeder(newUserRepo, newRoleRepo, newPermissionRepo)
	postSeeder := postDB.NewSeeder(cfg, newDisk, db, newSeedService, newApi, newPostRepo, newPostService, newCategoryRepo, categoryService)
	templateSeeder := templateDB.NewSeeder(cfg, newDisk, newTemplateRepo)
	infoBlockSeeder := infoBlockDB.NewSeeder(cfg, newDisk, newSeedService, newApi, newBlockService)
	messageSeeder := messageDB.NewMessageSeeder(newMessageService, newApi)

	seeder := migrate.NewSeeder(userSeeder, templateSeeder, infoBlockSeeder, postSeeder, messageSeeder, menuSeeder, settingsSeeder)

	newMigrator := migrate.NewMigrator(
		db.PostgreSQL(),
		userMigrator,
		postMigrator,
		infoBlockMigrator,
		galleryMigrator,
		templateMigrator,
		analyticMigrator,
		messageMigrator,
		fileMigrator,
		menuMigrator,
		settingsMigrator,
	)

	newScheduler := scheduler.NewScheduler(
		cfg,
		newQueue,
		fileProv,
	)

	newQueue.SetHandlers(map[string][]contract.QueueHandler{
		"messages": {
			messageQueue.NewMessageQueueHandler(newMessageService, newMessageCollectionService),
			mailerQueue.NewMailerQueueHandler(newMailer),
		},
		"users": {
			usersQueue.NewUserQueueHandler(newUserService),
		},
		"files": {
			fileQueue.NewFileQueueHandler(fileCollectionService),
		},
		"analytics": {
			analyticQueue.NewAnalyticQueueHandler(newAnalyticService, analyticCollectionService),
		},
		"posts": {
			menuQueue.NewPublisherQueueHandler(newMenuService, menuItemCollectionService),
		},
		"info_blocks": {
			postQueue.NewInfoBlockQueueHandler(newCategoriesService, newPostCollectionService, postTagCollectionService, newApi),
		},
		"galleries": {
			infoBlockQueue.NewGalleryQueueHandler(newBlockService, newBlockCollectionService, newBlockEventService),
			postQueue.NewGalleryQueueHandler(newCategoriesService, newPostCollectionService, postTagCollectionService, newApi),
		},
	})

	// Initialize Controllers
	adminWebPostController := postAdminWeb.NewWebPostController(
		newPostService,
		newPostCollectionService,
		categoryService,
		newCategoriesService,
		postTagCollectionService,
		newApi,
	)

	adminAjaxPostController := postAjax.NewPostController(
		newPostService,
		newPostCollectionService,
		categoryService,
		postTagCollectionService,
		newCategoriesService,
		newApi,
	)

	adminApiPostController := postApi.New(
		newPostService,
		categoryService,
		newCategoriesService,
		newApi,
	)

	adminWebCategoryController := postAdminWeb.NewWebCategoryController(
		newCategoriesService,
		categoryService,
		newApi,
	)

	adminAjaxCategoryController := postAjax.NewCategoryController(
		newCategoriesService,
		categoryService,
		newApi,
	)

	adminWebTagController := postAdminWeb.NewWebTagController(
		postTagService,
		postTagCollectionService,
		newApi,
	)

	adminAjaxTagController := postAjax.NewTagController(
		postTagService,
		postTagCollectionService,
		newApi,
	)

	adminWebInfoBlockController := infoBlockAdminWeb.NewInfoBlockWebController(
		newBlockService,
		newBlockCollectionService,
		newApi,
	)

	adminAjaxInfoBlockController := infoBlockAdminAjax.NewInfoBlockController(
		newBlockService,
		newBlockCollectionService,
		newApi,
	)

	adminWebTemplateController := templateAdminWeb.NewTemplateWebController(
		newTemplateService,
		newTemplateCollectionService,
		newApi,
	)

	adminAjaxTemplateController := templateAdminAjax.NewTemplateController(
		newTemplateService,
		newTemplateCollectionService,
		newApi,
	)

	adminWebMessageController := messageAdminWeb.NewMessageWebController(
		newMessageService,
		newMessageCollectionService,
		newApi,
	)

	adminAjaxMessageController := messageAdminAjax.NewMessageController(
		newMessageService,
		newMessageCollectionService,
		newApi,
	)

	adminAjaxGalleryController := galleryAjax.New(
		newGalleryRepo,
		newImageRepo,
		newImageService,
	)

	adminWebMenuController := menuAdminWeb.NewMenuWebController(
		newMenuService,
		newMenuCollectionService,
		menuItemService,
		menuItemCollectionService,
		newApi,
	)

	adminAjaxMenuController := menuAdminAjax.NewMenuAjaxController(
		newMenuService,
		newMenuCollectionService,
		newApi,
	)

	adminAjaxMenuItemController := menuAdminAjax.NewMenuItemAjaxController(
		menuItemService,
		menuItemCollectionService,
		newMenuService,
		newApi,
	)

	adminWebFileController := fileAdminWeb.NewFileController(
		uploadService,
		newFileService,
	)

	frontWebPostController := postFrontWeb.NewFrontWebController(
		cfg,
		newApi,
		newView,

		newPostService,
		newPostCollectionService,
		categoryService,
		newCategoriesService,
	)

	frontAjaxMessageController := messageFrontWeb.NewMessageController(
		newMailService,
	)

	frontWebUserController := userFrontWeb.NewUserWebController(
		newUserService,
		newAuthService,
		newCache,
	)

	adminAjaxPublisherController := publisherAjax.NewPublisherController(
		newPublisherService,
		newApi,
	)

	return &Container{
		Config:    cfg,
		Queue:     newQueue,
		Cache:     newCache,
		View:      newView,
		Scheduler: newScheduler,
		Disk:      newDisk,
		Minifier:  newMinifier,
		I18n:      newI18n,
		Store:     newStore,

		FileUploadService:     uploadService,
		FileCollectionService: fileCollectionService,
		FileService:           newFileService,
		FileProvider:          fileProv,

		GalleryResourceRepo: newResourceRepo,

		ImageRepo:     newImageRepo,
		ImageEvent:    newImageEvent,
		ImageService:  newImageService,
		ImageProvider: newImageProvider,

		GalleryProvider: newGalleryProvider,
		GalleryRepo:     newGalleryRepo,
		GalleryService:  newGalleryService,
		GalleryEvent:    newGalleryEvent,

		PostRepo:          newPostRepo,
		PostService:       newPostService,
		PostsService:      newPostCollectionService,
		PostProvider:      newBlogProvider,
		CategoryRepo:      newCategoryRepo,
		CategoriesService: newCategoriesService,
		CategoryService:   categoryService,

		TemplateProvider:          newTemplateProvider,
		TemplateRepo:              newTemplateRepo,
		TemplateService:           newTemplateService,
		TemplateCollectionService: newTemplateCollectionService,

		UserRepo:           newUserRepo,
		UserProvider:       newUserProvider,
		UserService:        newUserService,
		UserAuthService:    newAuthService,
		UserRoleRepo:       newRoleRepo,
		UserPermissionRepo: newPermissionRepo,

		AliasRepo:     newAliasRepo,
		AliasProvider: newAliasProvider,

		InfoBlockHasResourceRepo:   newInfoBlockHasResourceRepo,
		InfoBlockRepo:              newInfoBlockRepo,
		InfoBlockService:           newBlockService,
		InfoBlockCollectionService: newBlockCollectionService,
		InfoBlockProvider:          newBlockProvider,

		PostTagRepo:              postTagRepo,
		PostTagResourceRepo:      postTagResourceRepo,
		PostTagService:           postTagService,
		PostTagCollectionService: postTagCollectionService,

		MessageRepo:              newMessageRepo,
		MessageService:           newMessageService,
		MessageCollectionService: newMessageCollectionService,
		MailService:              newMailService,

		AnalyticRepo:              newAnalyticRepo,
		AnalyticService:           newAnalyticService,
		AnalyticCollectionService: analyticCollectionService,
		AnalyticProvider:          newAnalyticProvider,

		MenuRepo:                  menuRepo,
		MenuService:               newMenuService,
		MenuCollectionService:     newMenuCollectionService,
		MenuItemRepo:              menuItemRepo,
		MenuItemService:           menuItemService,
		MenuItemCollectionService: menuItemCollectionService,
		MenuProvider:              newMenuProvider,

		SettingsRepo:    newSettingsRepo,
		SettingsService: newSettingsService,

		Migrator: newMigrator,
		Seeder:   seeder,

		// Controllers
		AdminWebPostController:       adminWebPostController,
		AdminAjaxPostController:      adminAjaxPostController,
		AdminApiPostController:       adminApiPostController,
		AdminWebCategoryController:   adminWebCategoryController,
		AdminAjaxCategoryController:  adminAjaxCategoryController,
		AdminWebTagController:        adminWebTagController,
		AdminAjaxTagController:       adminAjaxTagController,
		AdminWebInfoBlockController:  adminWebInfoBlockController,
		AdminAjaxInfoBlockController: adminAjaxInfoBlockController,
		AdminWebTemplateController:   adminWebTemplateController,
		AdminAjaxTemplateController:  adminAjaxTemplateController,
		AdminWebMessageController:    adminWebMessageController,
		AdminAjaxMessageController:   adminAjaxMessageController,
		AdminAjaxGalleryController:   adminAjaxGalleryController,
		AdminWebMenuController:       adminWebMenuController,
		AdminAjaxMenuController:      adminAjaxMenuController,
		AdminAjaxMenuItemController:  adminAjaxMenuItemController,
		AdminWebFileController:       adminWebFileController,
		FrontWebPostController:       frontWebPostController,
		FrontAjaxMessageController:   frontAjaxMessageController,
		FrontWebUserController:       frontWebUserController,
		AdminAjaxPublisherController: adminAjaxPublisherController,

		Api: newApi,
	}
}
