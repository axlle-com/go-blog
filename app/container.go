package app

import (
	"context"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/cache"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/mailer"
	"github.com/axlle-com/blog/pkg/migrate"
	"github.com/axlle-com/blog/pkg/queue"
	"github.com/axlle-com/blog/pkg/view"

	analyticMigrate "github.com/axlle-com/blog/pkg/analytic/db/migrate"
	analyticProvider "github.com/axlle-com/blog/pkg/analytic/provider"
	analyticRepo "github.com/axlle-com/blog/pkg/analytic/repository"
	analyticService "github.com/axlle-com/blog/pkg/analytic/service"

	fileMigrate "github.com/axlle-com/blog/pkg/file/db/migrate"
	fileAdminWeb "github.com/axlle-com/blog/pkg/file/http"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
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
	infoBlockRepo "github.com/axlle-com/blog/pkg/info_block/repository"
	infoBlockService "github.com/axlle-com/blog/pkg/info_block/service"

	messageContracts "github.com/axlle-com/blog/pkg/message/contracts"
	messageDB "github.com/axlle-com/blog/pkg/message/db"
	messageMigrate "github.com/axlle-com/blog/pkg/message/db/migrate"
	messageAdminAjax "github.com/axlle-com/blog/pkg/message/http/admin/handlers/ajax"
	messageAdminWeb "github.com/axlle-com/blog/pkg/message/http/admin/handlers/web"
	messageFrontWeb "github.com/axlle-com/blog/pkg/message/http/front/handlers/ajax"
	messageRepo "github.com/axlle-com/blog/pkg/message/repository"
	messageService "github.com/axlle-com/blog/pkg/message/service"

	postDB "github.com/axlle-com/blog/pkg/post/db"
	postMigrate "github.com/axlle-com/blog/pkg/post/db/migrate"
	postAjax "github.com/axlle-com/blog/pkg/post/http/admin/handlers/ajax"
	postApi "github.com/axlle-com/blog/pkg/post/http/admin/handlers/api"
	postAdminWeb "github.com/axlle-com/blog/pkg/post/http/admin/handlers/web"
	postFrontWeb "github.com/axlle-com/blog/pkg/post/http/front/handlers/web"
	postRepo "github.com/axlle-com/blog/pkg/post/repository"
	postService "github.com/axlle-com/blog/pkg/post/service"

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
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
	usersService "github.com/axlle-com/blog/pkg/user/service"
)

type Container struct {
	Queue contracts.Queue
	Cache contracts.Cache
	View  contracts.View

	FileUploadService     *fileService.UploadService
	FileService           *fileService.Service
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
	PostsService      *postService.PostsService
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

	PostTagRepo         postRepo.PostTagRepository
	PostTagResourceRepo postRepo.PostTagResourceRepository
	PostTagService      *postService.PostTagService

	MessageRepo              messageContracts.MessageRepository
	MessageService           *messageService.MessageService
	MessageCollectionService *messageService.MessageCollectionService
	MailService              *messageService.MailService

	AnalyticRepo              analyticRepo.AnalyticRepository
	AnalyticService           *analyticService.AnalyticService
	AnalyticCollectionService *analyticService.AnalyticCollectionService
	AnalyticProvider          analyticProvider.AnalyticProvider

	Migrator contracts.Migrator
	Seeder   contracts.Seeder
}

func NewContainer(cfg contracts.Config, ctx context.Context) *Container {
	newQueue := queue.NewQueue()
	newQueue.StartWorkers(ctx, 4)

	newCache := cache.NewCache()
	newView := view.NewView(config.Config())

	mailerInterface := mailer.NewMailer(newQueue)

	newFileRepo := fileRepo.NewFileRepo()
	newFileService := fileService.NewService(newFileRepo)
	fileCollectionService := fileService.NewCollectionService(newFileRepo)
	uploadService := fileService.NewUploadService(newFileService)
	fileProv := fileProvider.NewProvider(uploadService, newFileService, fileCollectionService)

	newImageRepo := galleryRepo.NewImageRepo()
	newImageEvent := galleryService.NewImageEvent(fileProv)
	newImageService := galleryService.NewImageService(newImageRepo, newImageEvent, fileProv)
	newImageProvider := galleryProvider.NewImageProvider(newImageRepo)

	newResourceRepo := galleryRepo.NewResourceRepo()

	newGalleryRepo := galleryRepo.NewGalleryRepo()
	newGalleryEvent := galleryService.NewGalleryEvent(newImageService, newResourceRepo)
	newGalleryService := galleryService.NewGalleryService(newGalleryRepo, newGalleryEvent, newImageService, newResourceRepo, fileProv)
	newGalleryProvider := galleryProvider.NewProvider(newGalleryRepo, newGalleryService)

	newUserRepo := userRepository.NewUserRepo()
	newRoleRepo := userRepository.NewRoleRepo()
	newPermissionRepo := userRepository.NewPermissionRepo()
	newUserService := usersService.NewUserService(newUserRepo, newRoleRepo, newPermissionRepo)
	newAuthService := usersService.NewAuthService(newUserService)
	newUserProvider := userProvider.NewProvider(newUserRepo, newUserService)

	newTemplateRepo := templateRepo.NewTemplateRepo()
	newTemplateProvider := templateProvider.NewProvider(newTemplateRepo)
	newTemplateService := templateService.NewTemplateService(newTemplateRepo, newUserProvider)
	newTemplateCollectionService := templateService.NewTemplateCollectionService(newTemplateService, newTemplateRepo, newUserProvider)

	newMessageRepo := messageRepo.NewMessageRepo()
	newMessageService := messageService.NewMessageService(newMessageRepo, newUserProvider)
	newMessageCollectionService := messageService.NewMessageCollectionService(newMessageRepo, newMessageService, newUserProvider)
	newMailService := messageService.NewMailService(newMessageService, newMessageCollectionService, newUserProvider, mailerInterface, newQueue)

	newAliasRepo := alias.NewAliasRepo()
	newAliasProvider := alias.NewProvider(newAliasRepo)

	newPostRepo := postRepo.NewPostRepo()
	newCategoryRepo := postRepo.NewCategoryRepo()

	newInfoBlockHasResourceRepo := infoBlockRepo.NewResourceRepo()
	newInfoBlockRepo := infoBlockRepo.NewInfoBlockRepo()

	newBlockCollectionService := infoBlockService.NewInfoBlockCollectionService(newInfoBlockRepo, newInfoBlockHasResourceRepo, newGalleryProvider, newTemplateProvider, newUserProvider)
	newBlockService := infoBlockService.NewInfoBlockService(newInfoBlockRepo, newBlockCollectionService, newInfoBlockHasResourceRepo, newGalleryProvider, newTemplateProvider, newUserProvider, fileProv)
	newBlockProvider := infoBlockProvider.NewProvider(newBlockService, newBlockCollectionService)

	ptRepo := postRepo.NewPostTagRepo()
	ptrRepo := postRepo.NewResourceRepo()
	ptService := postService.NewPostTagService(ptRepo, ptrRepo)

	csService := postService.NewCategoriesService(newCategoryRepo, newAliasProvider, newGalleryProvider, newTemplateProvider, newUserProvider)
	cService := postService.NewCategoryService(newCategoryRepo, newAliasProvider, newGalleryProvider, fileProv, newBlockProvider)

	pService := postService.NewPostService(newPostRepo, csService, cService, newGalleryProvider, fileProv, newAliasProvider, newBlockProvider)
	psService := postService.NewPostsService(newPostRepo, csService, cService, newGalleryProvider, fileProv, newAliasProvider, newUserProvider, newTemplateProvider, newBlockProvider)

	newAnalyticRepo := analyticRepo.NewAnalyticRepo()
	newAnalyticService := analyticService.NewAnalyticService(newAnalyticRepo, newUserProvider)
	analyticCollectionService := analyticService.NewAnalyticCollectionService(newAnalyticRepo, newAnalyticService, newUserProvider)
	newAnalyticProvider := analyticProvider.NewAnalyticProvider(newAnalyticService, analyticCollectionService)

	mUser := userMigrate.NewMigrator(cfg.GetGORM())
	mPost := postMigrate.NewMigrator(cfg.GetGORM())
	mInfoBlock := infoBlockMigrate.NewMigrator(cfg.GetGORM())
	mGallery := galleryMigrate.NewMigrator(cfg.GetGORM())
	mTemplate := templateMigrate.NewMigrator(cfg.GetGORM())
	mAnalytic := analyticMigrate.NewMigrator(cfg.GetGORM())
	mMessage := messageMigrate.NewMigrator(cfg.GetGORM())
	mFile := fileMigrate.NewMigrator(cfg.GetGORM())

	sUser := userDB.NewSeeder(newUserRepo, newRoleRepo, newPermissionRepo)
	sPost := postDB.NewSeeder(newPostRepo, pService, newCategoryRepo, newUserProvider, newTemplateProvider)
	sTempl := templateDB.NewSeeder(newTemplateRepo)
	sInfo := infoBlockDB.NewSeeder(newBlockService, newTemplateProvider, newUserProvider)
	sMsg := messageDB.NewMessageSeeder(newMessageService, newUserProvider)

	seeder := migrate.NewSeeder(sUser, sTempl, sPost, sInfo, sMsg)

	newMigrator := migrate.NewMigrator(
		cfg.GetGORM(),
		mUser,
		mPost,
		mInfoBlock,
		mGallery,
		mTemplate,
		mAnalytic,
		mMessage,
		mFile,
	)

	return &Container{
		Queue: newQueue,
		Cache: newCache,
		View:  newView,

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
		PostService:       pService,
		PostsService:      psService,
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

		PostTagRepo:         ptRepo,
		PostTagResourceRepo: ptrRepo,
		PostTagService:      ptService,

		MessageRepo:              newMessageRepo,
		MessageService:           newMessageService,
		MessageCollectionService: newMessageCollectionService,
		MailService:              newMailService,

		AnalyticRepo:              newAnalyticRepo,
		AnalyticService:           newAnalyticService,
		AnalyticCollectionService: analyticCollectionService,
		AnalyticProvider:          newAnalyticProvider,

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

func (c *Container) PostController() postAjax.Controller {
	return postAjax.New(
		c.PostService,
		c.PostsService,
		c.CategoryService,
		c.CategoriesService,
		c.TemplateProvider,
		c.UserProvider,
		c.InfoBlockProvider,
	)
}

func (c *Container) PostWebController() postAdminWeb.Controller {
	return postAdminWeb.NewWebController(
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

func (c *Container) CategoryWebController() postAdminWeb.ControllerCategory {
	return postAdminWeb.NewWebControllerCategory(
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
