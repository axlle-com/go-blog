package routes

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	db2 "github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/cache"
	"github.com/axlle-com/blog/app/web"
	user "github.com/axlle-com/blog/pkg/user/models"
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

func StartWithLogin(ctx context.Context) (router *gin.Engine, cookie []*http.Cookie, err error) {
	router = SetupTestRouter(ctx)
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

func SetupTestRouter(ctx context.Context) *gin.Engine {
	if router == nil {
		cfg := config.Config()
		cfg.SetTestENV()

		db2.InitDB(cfg)

		container := app.NewContainer(cfg, ctx)

		err := container.Migrator.Migrate()
		if err != nil {
			return nil
		}

		err = container.Seeder.Seed()
		if err != nil {
			return nil
		}

		gob.Register(user.User{})

		router = gin.New()

		store := models.Store(cfg)
		router.Use(sessions.Sessions(cfg.SessionsName(), store))

		web.NewTemplate(router)
		InitializeApiRoutes(router, container)
		InitializeWebRoutes(router, container)
		cache.NewCache().ResetUsersSession()
	}
	return router
}
