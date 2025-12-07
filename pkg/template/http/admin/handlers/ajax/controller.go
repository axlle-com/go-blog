package ajax

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/service"
	"github.com/gin-gonic/gin"
)

type TemplateController interface {
	GetTemplate(ctx *gin.Context)
	GetResourceTemplate(ctx *gin.Context)
	UpdateTemplate(*gin.Context)
	CreateTemplate(*gin.Context)
	DeleteTemplate(*gin.Context)
	FilterTemplate(*gin.Context)
}

func NewTemplateController(
	templateService *service.TemplateService,
	templateCollectionService *service.TemplateCollectionService,
	api *api.Api,
) TemplateController {
	return &templateController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		api:                       api,
	}
}

type templateController struct {
	*app.BaseAjax

	templateService           *service.TemplateService
	templateCollectionService *service.TemplateCollectionService
	api                       *api.Api
}

func ShowIndexPageTest(ctx *gin.Context) {
	fileName := filepath.Base("index.gohtml")
	templatePath := filepath.Join("src/templates", fileName)
	data, err := os.ReadFile(templatePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка чтения файла: %s", err.Error())
		return
	}

	var base app.BaseAjax
	tFunc := base.BuildT(ctx)
	templateData := app.PrepareTemplateData(ctx, gin.H{
		"title":   "Home Page",
		"payload": string(data),
	}, tFunc)
	ctx.HTML(http.StatusOK, "test", templateData)
}
