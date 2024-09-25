package main

import (
	"encoding/gob"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	gob.Register(user.User{})
	cfg := config.GetConfig()
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := db.InitRedis(cfg)
	router.Use(sessions.Sessions(config.SessionsName, store))

	db.Init(cfg.DBUrl)

	//web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeApiRoutes(router)
	routes.InitializeWebRoutes(router)

	err = router.Run(cfg.Port)
	if err != nil {
		panic("Error run")
	}
}
