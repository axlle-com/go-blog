package routes

import (
	"bytes"
	"encoding/gob"
	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	db2 "github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/web"
	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	mTemplate "github.com/axlle-com/blog/pkg/template/db/migrate"
	dbUser "github.com/axlle-com/blog/pkg/user/db"
	mUser "github.com/axlle-com/blog/pkg/user/db/migrate"
	user "github.com/axlle-com/blog/pkg/user/models"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

var router *gin.Engine

func PerformLogin(router *gin.Engine) ([]*http.Cookie, error) {
	requestBody := `{"email":"admin@admin.ru","password":"123456"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code >= 400 {
		return nil, http.ErrNoCookie
	}

	cookie := w.Result().Cookies()
	return cookie, nil
}

func StartWithLogin() (router *gin.Engine, cookie []*http.Cookie, err error) {
	router = SetupTestRouter()
	requestBody := `{"email":"admin@admin.ru","password":"123456"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code >= 400 {
		return nil, nil, http.ErrNoCookie
	}

	cookie = w.Result().Cookies()
	return router, cookie, nil
}

func SetupTestRouter() *gin.Engine {
	if router == nil {
		cfg := config.Config()
		cfg.SetTestENV()

		db2.Init(cfg.DBUrlTest())

		container := app.New()

		mUser.NewMigrator().Migrate()
		mTemplate.NewMigrator().Migrate()
		mPost.NewMigrator().Migrate()
		mGallery.NewMigrator().Migrate()

		dbUser.NewSeeder(
			container.UserRepo,
			userRepository.NewRoleRepo(),
			userRepository.NewPermissionRepo(),
		).Seed()

		gob.Register(user.User{})

		router = gin.New()

		store := models.Store(cfg.RedisHost(), "", cfg.KeyCookie())
		router.Use(sessions.Sessions(cfg.SessionsName(), store))

		web.NewTemplate(router)
		InitializeApiRoutes(router, container)
		InitializeWebRoutes(router, container)
		db2.NewCache().ResetUsersSession()
	}
	return router
}
