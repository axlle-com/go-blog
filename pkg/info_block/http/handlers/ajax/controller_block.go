package ajax

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type InfoBlockController interface {
	UpdateInfoBlock(*gin.Context)
	CreateInfoBlock(*gin.Context)
	DeleteInfoBlock(*gin.Context)
	FilterInfoBlock(*gin.Context)
}

func NewInfoBlockController(
	blockService *service.InfoBlockService,
	blockCollectionService *service.InfoBlockCollectionService,
	template template.TemplateProvider,
	user user.UserProvider,
) InfoBlockController {
	return &blockController{
		blockService:           blockService,
		blockCollectionService: blockCollectionService,
		templateProvider:       template,
		userProvider:           user,
	}
}

type blockController struct {
	*app.BaseAjax

	blockService           *service.InfoBlockService
	blockCollectionService *service.InfoBlockCollectionService
	templateProvider       template.TemplateProvider
	userProvider           user.UserProvider
}
