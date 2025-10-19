package app

import (
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/cache"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/service/mailer"
	mailerQueue "github.com/axlle-com/blog/app/service/mailer/queue"
	"github.com/axlle-com/blog/app/service/migrate"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/app/service/scheduler"
	"github.com/axlle-com/blog/app/service/storage"
	"github.com/axlle-com/blog/app/service/view"

	"github.com/axlle-com/blog/pkg/blog/provider"

	"github.com/axlle-com/blog/pkg/alias"
	analyticMigrate "github.com/axlle-com/blog/pkg/analytic/db/migrate"
	analyticProvider "github.com/axlle-com/blog/pkg/analytic/provider"
	analyticQueue "github.com/axlle-com/blog/pkg/analytic/queue"
	analyticRepo "github.com/axlle-com/blog/pkg/analytic/repository"
	analyticService "github.com/axlle-com/blog/pkg/analytic/service"

	fileMigrate "github.com/axlle-com/blog/pkg/file/db/migrate"
	fileAdminWeb "github.com/axlle-com/blog/pkg/file/http"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	fileQueue "github.com/axlle-com/blog/pkg/file/queue"
	fileRepo "github.com/axlle-com/blog/pkg/file/repository"
	fileService "github.com/axlle-com/blog/pkg/file/service"

	menuDB "github.com/axlle-com/blog/pkg/menu/db"
	menuMigrate "github.com/axlle-com/blog/pkg/menu/db/migrate"
	menuAdminAjax "github.com/axlle-com/blog/pkg/menu/http/handlers/ajax"
	menuAdminWeb "github.com/axlle-com/blog/pkg/menu/http/handlers/web"
	menuQueue "github.com/axlle-com/blog/pkg/menu/queue"
	menuRepository "github.com/axlle-com/blog/pkg/menu/repository"
	menuService "github.com/axlle-com/blog/pkg/menu/service"

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
	infoBlockRepo "github.com/axlle-com/blog/pkg/info_block/repository"
	infoBlockService "github.com/axlle-com/blog/pkg/info_block/service"

	messageContracts "github.com/axlle-com/blog/pkg/message/contracts"
	messageDB "github.com/axlle-com/blog/pkg/message/db"
	messageMigrate "github.com/axlle-com/blog/pkg/message/db/migrate"
	messageAdminAjax "github.com/axlle-com/blog/pkg/message/http/admin/handlers/ajax"
	messageAdminWeb "github.com/axlle-com/blog/pkg/message/http/admin/handlers/web"
	messageFrontWeb "github.com/axlle-com/blog/pkg/message/http/front/handlers/ajax"
	messageQueue "github.com/axlle-com/blog/pkg/message/queue"
	messageRepo "github.com/axlle-com/blog/pkg/message/repository"
	messageService "github.com/axlle-com/blog/pkg/message/service"

	postDB "github.com/axlle-com/blog/pkg/blog/db"
	postMigrate "github.com/axlle-com/blog/pkg/blog/db/migrate"
	postAjax "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/api"
	postAdminWeb "github.com/axlle-com/blog/pkg/blog/http/admin/handlers/web"
	postFrontWeb "github.com/axlle-com/blog/pkg/blog/http/front/handlers/web"
	postRepo "github.com/axlle-com/blog/pkg/blog/repository"
	postService "github.com/axlle-com/blog/pkg/blog/service"

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
)

type Container struct {
	Config    contracts.Config
	Queue     contracts.Queue
	Cache     contracts.Cache
	View      contracts.View
	Scheduler contracts.Scheduler

	FileUploadService     *fileService.UploadService
	FileService           *fileService.FileService
	FileCollectionService *fileService.CollectionService
	FileProvider          fileProvider.FileProvider

	ImageRepo     galleryRepo.GalleryImageRepository
	ImageEvent    *galleryService.ImageEvent
	ImageService  *galleryService.ImageService
	ImageProvider galleryProvider.ImageProvider

	GalleryRepo         galleryRepo.GalleryRepository
	GalleryEvent        *galleryService.GalleryEvent
	GalleryService      *galleryService.GalleryService
	GalleryProvider     galleryProvider.GalleryProvider
	GalleryResourceRepo galleryRepo.GalleryResourceRepository

	PostRepo          postRepo.PostRepository
	PostService       *postService.PostService
	PostsService      *postService.PostCollectionService
	PostProvider      contracts.PostProvider
	CategoryRepo      postRepo.CategoryRepository
	CategoriesService *postService.CategoriesService
	CategoryService   *postService.CategoryService

	TemplateProvider          templateProvider.TemplateProvider
	TemplateRepo              templateRepo.TemplateRepository
	TemplateService           *templateService.TemplateService
	TemplateCollectionService *templateService.TemplateCollectionService

	UserRepo        userRepository.UserRepository
	UserProvider    userProvider.UserProvider
	UserService     *usersService.UserService
	UserAuthService *usersService.AuthService

	UserRoleRepo       userRepository.RoleRepository
	UserPermissionRepo userRepository.PermissionRepository

	AliasRepo     alias.AliasRepository
	AliasProvider alias.AliasProvider

	InfoBlockHasResourceRepo   infoBlockRepo.InfoBlockHasResourceRepository
	InfoBlockRepo              infoBlockRepo.InfoBlockRepository
	InfoBlockService           *infoBlockService.InfoBlockService
	InfoBlockCollectionService *infoBlockService.InfoBlockCollectionService
	InfoBlockProvider          infoBlockProvider.InfoBlockProvider

	PostTagRepo              postRepo.PostTagRepository
	PostTagResourceRepo      postRepo.PostTagResourceRepository
	PostTagService           *postService.TagService
	PostTagCollectionService *postService.TagCollectionService

	MessageRepo              messageContracts.MessageRepository
	MessageService           *messageService.MessageService
	MessageCollectionService *messageService.MessageCollectionService
	MailService              *messageService.MailService

	AnalyticRepo              analyticRepo.AnalyticRepository
	AnalyticService           *analyticService.AnalyticService
	AnalyticCollectionService *analyticService.AnalyticCollectionService
	AnalyticProvider          analyticProvider.AnalyticProvider

	MenuRepo                  menuRepository.MenuRepository
	MenuService               *menuService.MenuService
	MenuCollectionService     *menuService.MenuCollectionService
	MenuItemRepo              menuRepository.MenuItemRepository
	MenuItemService           *menuService.MenuItemService
	MenuItemCollectionService *menuService.MenuItemCollectionService

	Migrator contracts.Migrator
	Seeder   contracts.Seeder
}

func NewContainer(cfg contracts.Config, db contracts.DB) *Container {
	newQueue := queue.NewQueue()
	newCache := cache.NewCache()
	newView := view.NewView(config.Config())
	newMailer := mailer.NewMailer(cfg, newQueue)

	newFileRepo := fileRepo.NewFileRepo(db.PostgreSQL())
	newFileService := fileService.NewFileService(newFileRepo)
	newStorageService := storage.NewLocalStorageService(cfg)
	uploadService := fileService.NewUploadService(newFileService, newStorageService)
	fileCollectionService := fileService.NewCollectionService(newFileRepo, newFileService, uploadService)
	fileProv := fileProvider.NewFileProvider(uploadService, newFileService, fileCollectionService)

	newImageRepo := galleryRepo.NewImageRepo(db.PostgreSQL())
	newImageEvent := galleryService.NewImageEvent(fileProv)
	newImageService := galleryService.NewImageService(newImageRepo, newImageEvent, fileProv)
	newImageProvider := galleryProvider.NewImageProvider(newImageRepo)

	newResourceRepo := galleryRepo.NewResourceRepo(db.PostgreSQL())

	newGalleryRepo := galleryRepo.NewGalleryRepo(db.PostgreSQL())
	newGalleryEvent := galleryService.NewGalleryEvent(newImageService, newResourceRepo)
	newGalleryService := galleryService.NewGalleryService(newGalleryRepo, newGalleryEvent, newImageService, newResourceRepo, fileProv)
	newGalleryProvider := galleryProvider.NewProvider(newGalleryRepo, newGalleryService)

	newUserRepo := userRepository.NewUserRepo(db.PostgreSQL())
	newRoleRepo := userRepository.NewRoleRepo(db.PostgreSQL())
	newPermissionRepo := userRepository.NewPermissionRepo(db.PostgreSQL())
	newUserService := usersService.NewUserService(newUserRepo, newRoleRepo, newPermissionRepo)
	newAuthService := usersService.NewAuthService(newUserService)
	newUserProvider := userProvider.NewProvider(newUserRepo, newUserService)

	newTemplateRepo := templateRepo.NewTemplateRepo(db.PostgreSQL())
	newTemplateProvider := templateProvider.NewProvider(newTemplateRepo)
	newTemplateService := templateService.NewTemplateService(newTemplateRepo, newUserProvider)
	newTemplateCollectionService := templateService.NewTemplateCollectionService(newTemplateService, newTemplateRepo, newUserProvider)

	newMessageRepo := messageRepo.NewMessageRepo(db.PostgreSQL())
	newMessageService := messageService.NewMessageService(newMessageRepo, newUserProvider)
	newMessageCollectionService := messageService.NewMessageCollectionService(newMessageRepo, newMessageService, newUserProvider)
	newMailService := messageService.NewMailService(cfg, newQueue, newMessageService, newMessageCollectionService, newUserProvider)

	newAliasRepo := alias.NewAliasRepo(db.PostgreSQL())
	newAliasProvider := alias.NewAliasProvider(newAliasRepo)

	newPostRepo := postRepo.NewPostRepo(db.PostgreSQL())
	newCategoryRepo := postRepo.NewCategoryRepo(db.PostgreSQL())

	newInfoBlockHasResourceRepo := infoBlockRepo.NewResourceRepo(db.PostgreSQL())
	newInfoBlockRepo := infoBlockRepo.NewInfoBlockRepo(db.PostgreSQL())

	newBlockCollectionService := infoBlockService.NewInfoBlockCollectionService(newInfoBlockRepo, newInfoBlockHasResourceRepo, newGalleryProvider, newTemplateProvider, newUserProvider)
	newBlockService := infoBlockService.NewInfoBlockService(newInfoBlockRepo, newBlockCollectionService, newInfoBlockHasResourceRepo, newGalleryProvider, newTemplateProvider, newUserProvider, fileProv)
	newBlockProvider := infoBlockProvider.NewProvider(newBlockService, newBlockCollectionService)

	postTagRepo := postRepo.NewPostTagRepo(db.PostgreSQL())
	postTagResourceRepo := postRepo.NewResourceRepo(db.PostgreSQL())
	postTagService := postService.NewTagService(postTagRepo, postTagResourceRepo, newAliasProvider, newGalleryProvider, newBlockProvider, fileProv)
	postTagCollectionService := postService.NewTagCollectionService(postTagService, postTagRepo, postTagResourceRepo, newTemplateProvider)

	csService := postService.NewCategoriesService(newCategoryRepo, newAliasProvider, newGalleryProvider, newTemplateProvider, newUserProvider)
	cService := postService.NewCategoryService(newCategoryRepo, newAliasProvider, newGalleryProvider, fileProv, newBlockProvider)

	newPostService := postService.NewPostService(newQueue, newPostRepo, csService, cService, postTagCollectionService, newGalleryProvider, fileProv, newAliasProvider, newBlockProvider)
	newPostCollectionService := postService.NewPostCollectionService(newPostRepo, csService, cService, newGalleryProvider, fileProv, newAliasProvider, newUserProvider, newTemplateProvider, newBlockProvider)
	newPostProvider := provider.NewPostProvider(newPostService, newPostCollectionService, csService, postTagCollectionService)
	newAnalyticRepo := analyticRepo.NewAnalyticRepo(db.PostgreSQL())
	newAnalyticService := analyticService.NewAnalyticService(newAnalyticRepo, newUserProvider)
	analyticCollectionService := analyticService.NewAnalyticCollectionService(newAnalyticRepo, newAnalyticService, newUserProvider)
	newAnalyticProvider := analyticProvider.NewAnalyticProvider(newAnalyticService, analyticCollectionService)

	menuRepo := menuRepository.NewMenuRepo(db.PostgreSQL())
	menuItemRepo := menuRepository.NewMenuItemRepo(db.PostgreSQL())
	menuItemService := menuService.NewMenuItemService(menuItemRepo)
	menuItemCollectionService := menuService.NewMenuItemCollectionService(menuItemRepo, menuItemService)
	newMenuService := menuService.NewMenuService(menuRepo, menuItemService, menuItemCollectionService)
	newMenuCollectionService := menuService.NewMenuCollectionService(menuRepo, newMenuService)

	menuSeeder := menuDB.NewMenuSeeder(menuRepo, menuItemRepo, newPostProvider, newTemplateProvider)

	userMigrator := userMigrate.NewMigrator(db.PostgreSQL())
	postMigrator := postMigrate.NewMigrator(db.PostgreSQL())
	infoBlockMigrator := infoBlockMigrate.NewMigrator(db.PostgreSQL())
	galleryMigrator := galleryMigrate.NewMigrator(db.PostgreSQL())
	templateMigrator := templateMigrate.NewMigrator(db.PostgreSQL())
	analyticMigrator := analyticMigrate.NewMigrator(db.PostgreSQL())
	messageMigrator := messageMigrate.NewMigrator(db.PostgreSQL())
	fileMigrator := fileMigrate.NewMigrator(db.PostgreSQL())
	menuMigrator := menuMigrate.NewMigrator(db.PostgreSQL())

	userSeeder := userDB.NewSeeder(newUserRepo, newRoleRepo, newPermissionRepo)
	postSeeder := postDB.NewSeeder(newPostRepo, newPostService, newCategoryRepo, newUserProvider, newTemplateProvider)
	templateSeeder := templateDB.NewSeeder(newTemplateRepo)
	infoBlockSeeder := infoBlockDB.NewSeeder(newBlockService, newTemplateProvider, newUserProvider)
	messageSeeder := messageDB.NewMessageSeeder(newMessageService, newUserProvider)

	seeder := migrate.NewSeeder(userSeeder, templateSeeder, postSeeder, infoBlockSeeder, messageSeeder, menuSeeder)

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
	)

	newScheduler := scheduler.NewScheduler(
		cfg,
		newQueue,
		fileProv,
	)

	newQueue.SetHandlers(map[string][]contracts.QueueHandler{
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
	})

	return &Container{
		Config:    cfg,
		Queue:     newQueue,
		Cache:     newCache,
		View:      newView,
		Scheduler: newScheduler,

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
		PostProvider:      newPostProvider,
		CategoryRepo:      newCategoryRepo,
		CategoriesService: csService,
		CategoryService:   cService,

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

		Migrator: newMigrator,
		Seeder:   seeder,
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

func (c *Container) PostController() postAjax.PostController {
	return postAjax.NewPostController(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.PostTagCollectionService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) PostWebController() postAdminWeb.PostController {
	return postAdminWeb.NewWebPostController(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.PostTagCollectionService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) PostTagWebController() postAdminWeb.TagController {
	return postAdminWeb.NewWebTagController(
		c.PostTagService,
		c.PostTagCollectionService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) PostTagAjaxController() postAjax.TagController {
	return postAjax.NewTagController(
		c.PostTagService,
		c.PostTagCollectionService,
		c.TemplateProvider,
		c.UserProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) CategoryWebController() postAdminWeb.CategoryController {
	return postAdminWeb.NewWebCategoryController(
		c.CategoriesService,
		c.CategoryService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) CategoryController() postAjax.CategoryController {
	return postAjax.NewCategoryController(
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

func (c *Container) InfoBlockController() infoBlockAdminAjax.InfoBlockController {
	return infoBlockAdminAjax.NewInfoBlockController(
		c.InfoBlockService,
		c.InfoBlockCollectionService,
		c.TemplateProvider,
		c.UserProvider,
	)
}

func (c *Container) InfoBlockWebController() infoBlockAdminWeb.InfoBlockWebController {
	return infoBlockAdminWeb.NewInfoBlockWebController(
		c.InfoBlockService,
		c.InfoBlockCollectionService,
		c.TemplateProvider,
		c.UserProvider,
		c.GalleryProvider,
	)
}

func (c *Container) PostFrontWebController() postFrontWeb.PostController {
	return postFrontWeb.NewFrontWebController(
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

func (c *Container) TemplateWebController() templateAdminWeb.TemplateWebController {
	return templateAdminWeb.NewTemplateWebController(
		c.TemplateService,
		c.TemplateCollectionService,
		c.UserProvider,
	)
}

func (c *Container) TemplateController() templateAdminAjax.TemplateController {
	return templateAdminAjax.NewTemplateController(
		c.TemplateService,
		c.TemplateCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageController() messageAdminWeb.MessageWebController {
	return messageAdminWeb.NewMessageWebController(
		c.MessageService,
		c.MessageCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageAjaxController() messageAdminAjax.MessageController {
	return messageAdminAjax.NewMessageController(
		c.MessageService,
		c.MessageCollectionService,
		c.UserProvider,
	)
}

func (c *Container) MessageFrontController() messageFrontWeb.MessageController {
	return messageFrontWeb.NewMessageController(
		c.MailService,
	)
}

func (c *Container) UserFrontController() userFrontWeb.Controller {
	return userFrontWeb.NewUserWebController(
		c.UserService,
		c.UserAuthService,
		c.Cache,
	)
}

func (c *Container) FileController() fileAdminWeb.Controller {
	return fileAdminWeb.NewFileController(
		c.FileUploadService,
		c.FileService,
	)
}

func (c *Container) MenuController() menuAdminWeb.Controller {
	return menuAdminWeb.NewMenuWebController(
		c.MenuService,
		c.MenuCollectionService,
		c.MenuItemService,
		c.MenuItemCollectionService,
		c.TemplateProvider,
		c.PostProvider,
	)
}

func (c *Container) MenuItemAjaxController() menuAdminAjax.ControllerMenuItem {
	return menuAdminAjax.NewMenuItemAjaxController(
		c.MenuItemService,
		c.MenuItemCollectionService,
	)
}

func (c *Container) MenuAjaxController() menuAdminAjax.ControllerMenu {
	return menuAdminAjax.NewMenuAjaxController(
		c.MenuService,
		c.MenuCollectionService,
		c.TemplateProvider,
		c.PostProvider,
	)
}
