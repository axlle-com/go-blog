package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/axlle-com/blog/pkg/app"
	"github.com/axlle-com/blog/pkg/app/middleware"
	"github.com/axlle-com/blog/pkg/app/web"
	file "github.com/axlle-com/blog/pkg/file/http"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
)

func InitializeWebRoutes(r *gin.Engine, container *app.Container) {
	postFrontWebController := container.PostFrontWebController()
	postController := container.PostController()
	postWebController := container.PostWebController()
	postCategoryWebController := container.CategoryWebController()
	postCategoryController := container.CategoryController()
	galleryController := container.GalleryAjaxController()

	fileController := file.NewFileController(
		container.FileService,
	)

	userController := user.NewUserWebController(
		container.UserService,
		container.UserAuthService,
	)

	infoBlockController := container.InfoBlockWebController()
	infoBlockAjaxController := container.InfoBlockController()

	r.Use(middleware.Error())
	r.Use(middleware.Analytic())
	r.GET("/", postFrontWebController.GetHome)
	r.GET("/test", ShowIndexPageTest)
	r.POST("/test", SavePageTest)
	r.GET("/login", userController.Login)
	r.POST("/auth", userController.Auth)
	r.POST("/user", userController.CreateUser)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("", userController.Index)
		protected.GET("/logout", userController.Logout)

		protected.POST("/file/image", fileController.UploadImage)
		protected.POST("/file/images", fileController.UploadImages)
		protected.DELETE("/file/image/*filePath", fileController.DeleteImage)

		protected.GET("/posts", postWebController.GetPosts)
		protected.GET("/posts/:id", postWebController.GetPost)
		protected.GET("/posts/form", postWebController.CreatePost)
		protected.POST("/posts", postController.CreatePost)
		protected.GET("/posts/filter", postController.FilterPosts)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)
		protected.DELETE("/posts/:id/image", postController.DeletePostImage)

		protected.GET("/categories", postCategoryWebController.GetCategories)
		protected.GET("/categories/:id", postCategoryWebController.GetCategory)
		protected.POST("/categories", postCategoryController.CreateCategory)
		protected.PUT("/categories/:id", postCategoryController.UpdateCategory)
		protected.DELETE("/categories/:id/image", postCategoryController.DeleteCategoryImage)
		protected.DELETE("/categories/:id", postCategoryController.DeleteCategory)
		protected.GET("/categories/form", postCategoryWebController.CreateCategory)
		protected.GET("/categories/filter", postCategoryController.FilterCategory)

		protected.GET("/info-blocks", infoBlockController.GetInfoBlocks)
		protected.GET("/info-blocks/:id", infoBlockController.GetInfoBlock)

		protected.POST("/info-blocks", infoBlockAjaxController.CreateInfoBlock)
		protected.PUT("/info-blocks/:id", infoBlockAjaxController.UpdateInfoBlock)
		protected.DELETE("/info-blocks/:id", infoBlockAjaxController.DeleteInfoBlock)
		protected.GET("/info-blocks/filter", infoBlockAjaxController.FilterInfoBlock)

		protected.DELETE("/gallery/:id/image/:image_id", galleryController.DeleteImage)
	}
	r.GET("/:alias", postFrontWebController.GetPost)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/admin") {
			ctx.HTML(http.StatusNotFound, "admin.404", gin.H{
				"title": "Админка — 404",
				"menu":  menu.NewMenu(ctx.FullPath()),
			})
		} else {
			ctx.HTML(http.StatusNotFound, "404", gin.H{
				"title": "Страница не найдена",
			})
		}
	})
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
