package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/web"
	"github.com/axlle-com/blog/pkg/template/service"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
	userProvider userProvider.UserProvider,
) TemplateController {
	return &templateController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		userProvider:              userProvider,
	}
}

type templateController struct {
	*app.BaseAjax

	templateService           *service.TemplateService
	templateCollectionService *service.TemplateCollectionService
	userProvider              userProvider.UserProvider
}

func ShowIndexPageTest(c *gin.Context) {
	fileName := filepath.Base("index.gohtml")
	templatePath := filepath.Join("src/templates", fileName)
	data, err := os.ReadFile(templatePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка чтения файла: %s", err.Error())
		return
	}

	c.HTML(
		http.StatusOK,
		"test",
		gin.H{
			"title":   "Home Page",
			"payload": string(data),
		},
	)
}

func SavePageTest(c *gin.Context) {
	code := c.PostForm("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Не передано содержимое шаблона (code)")
		return
	}

	fileName := filepath.Base("index.gohtml")
	templatePath := filepath.Join("src/templates", fileName)

	err := os.WriteFile(templatePath, []byte(code), 0644)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка записи файла: %s", err.Error())
		return
	}
	web.NewTemplate(nil).ReLoad()
	c.String(http.StatusOK, "Файл успешно сохранён")
}

func SavePageTest2(c *gin.Context) {
	code := c.PostForm("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Не передано содержимое шаблона (code)")
		return
	}

	err := web.NewTemplate(nil).AddTemplateFromString("index", code)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, "Файл успешно сохранён")
}
