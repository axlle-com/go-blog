package main

import (
	"encoding/gob"
	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	routes2 "github.com/axlle-com/blog/app/routes"
	"github.com/axlle-com/blog/app/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	user "github.com/axlle-com/blog/pkg/user/models"
)

func main() {
	cfg := config.Config()
	router := Init(cfg)
	err := router.Run(cfg.Port())
	if err != nil {
		panic("Error run")
	}
}

func Init(cfg contracts.Config) *gin.Engine {
	gob.Register(user.User{})

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := models.Store(cfg.RedisHost(), "", cfg.KeyCookie())
	router.Use(sessions.Sessions(cfg.SessionsName(), store))

	db.Init(cfg.DBUrl())

	container := app.New()

	web.InitMinify()
	web.NewTemplate(router)
	routes2.InitializeApiRoutes(router, container)
	routes2.InitializeWebRoutes(router, container)
	return router
}
