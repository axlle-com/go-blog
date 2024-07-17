package main

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := cookie.NewStore([]byte(cfg.KeyCookie))
	store.Options(sessions.Options{
		MaxAge: 86400 * 7, // 7 дней
		Path:   "/",
	})
	router.Use(sessions.Sessions(config.SessionsName, store))

	web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeWebRoutes(router)
	routes.InitializeApiRoutes(router)

	err = router.Run(cfg.Port)
	if err != nil {
		panic("Error run")
	}
}
