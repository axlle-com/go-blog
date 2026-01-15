package routes

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/http/httptest"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/middleware"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/service/cache"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

		newDB, err := db.SetupDB(cfg)
		if err != nil {
			panic("db not initialized")
		}

		container := di.NewContainer(cfg, newDB)

		err = container.Migrator.Migrate()
		if err != nil {
			return nil
		}

		err = container.Seeder.Seed()
		if err != nil {
			return nil
		}

		gob.Register(user.User{})

		router = gin.New()

		languageMiddleware := middleware.NewLanguage(container.I18n)

		store := models.Store(cfg)
		router.Use(sessions.Sessions(cfg.SessionsName(), store))
		router.Use(languageMiddleware.Handler())

		container.View.SetRouter(router)
		container.View.Load()
		container.Disk.SetupStaticFiles(router)

		InitApiRoutes(router, container)
		InitWebRoutes(router, cfg, container)
		cache.NewCache().ResetUsersSession()
	}
	return router
}
