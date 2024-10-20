package main

import (
	"encoding/gob"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

	store := db.RedisStore(cfg.RedisHost(), "", cfg.KeyCookie())
	router.Use(sessions.Sessions(cfg.SessionsName(), store))

	db.Init(cfg.DBUrl())

	//web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeApiRoutes(router)
	routes.InitializeWebRoutes(router)
	return router
}
