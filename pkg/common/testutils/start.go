package testutils

import (
	"bytes"
	"encoding/gob"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	mTemplate "github.com/axlle-com/blog/pkg/template/db/migrate"
	dbUser "github.com/axlle-com/blog/pkg/user/db"
	mUser "github.com/axlle-com/blog/pkg/user/db/migrate"
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

func SetupTestRouter() *gin.Engine {
	if router == nil {
		mUser.Migrate()
		mTemplate.Migrate()
		mPost.Migrate()
		mGallery.Migrate()

		dbUser.SeedPermissions()
		dbUser.SeedRoles()
		dbUser.SeedUsersDefault()
		gob.Register(user.User{})
		router = gin.New()

		cfg := config.Config()
		cfg.SetTestENV()
		store := db.RedisStore(cfg.RedisHost(), "", cfg.KeyCookie())
		router.Use(sessions.Sessions(cfg.SessionsName(), store))

		db.Init(cfg.DBUrlTest())

		web.InitTemplate(router)
		routes.InitializeApiRoutes(router)
		routes.InitializeWebRoutes(router)
		db.Cache().ResetUsersSession()
	}
	return router
}
