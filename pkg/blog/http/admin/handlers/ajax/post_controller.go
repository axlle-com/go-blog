package ajax

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
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
	templateProvider template.TemplateProvider,
	user user.UserProvider,
	infoBlock appPovider.InfoBlockProvider,

) PostController {
	return &postController{
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		tagCollectionService:  tagCollectionService,
		templateProvider:      templateProvider,
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
	templateProvider      template.TemplateProvider
	user                  user.UserProvider
	infoBlock             appPovider.InfoBlockProvider
}

func (c *postController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.templateProvider.GetForResources(&models.Post{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
