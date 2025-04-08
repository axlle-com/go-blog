package ajax

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/info_block/provider"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	AddPostInfoBlock(*gin.Context)
	DeletePost(*gin.Context)
	DeletePostImage(*gin.Context)
	FilterPosts(*gin.Context)
}

func New(
	service *service.PostService,
	services *service.PostsService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	infoBlock provider.InfoBlockProvider,

) Controller {
	return &controller{
		postService:       service,
		postsService:      services,
		categoryService:   category,
		categoriesService: categories,
		template:          template,
		user:              user,
		infoBlock:         infoBlock,
	}
}

type controller struct {
	*app.BaseAjax

	postService       *service.PostService
	postsService      *service.PostsService
	categoryService   *service.CategoryService
	categoriesService *service.CategoriesService
	template          template.TemplateProvider
	user              user.UserProvider
	infoBlock         provider.InfoBlockProvider
}
