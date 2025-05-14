package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/axlle-com/blog/pkg/info_block/provider"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	AddPostInfoBlock(*gin.Context)
	DeletePost(*gin.Context)
	DeletePostImage(*gin.Context)
	FilterPosts(*gin.Context)
}

func NewPostController(
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	tagCollectionService *service.TagCollectionService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	infoBlock provider.InfoBlockProvider,

) PostController {
	return &postController{
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		tagCollectionService:  tagCollectionService,
		template:              template,
		user:                  user,
		infoBlock:             infoBlock,
	}
}

type postController struct {
	*app.BaseAjax

	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	tagCollectionService  *service.TagCollectionService
	template              template.TemplateProvider
	user                  user.UserProvider
	infoBlock             provider.InfoBlockProvider
}
