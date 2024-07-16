package main

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic("Ошибка загрузки конфигурации: " + err.Error())
	}

	cfg := config.GetConfig()
	router := gin.Default()

	h := db.Init(cfg.DBUrl)

	web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeWebRoutes(router, h)
	routes.InitializeApiRoutes(router, h)

	err = router.Run(cfg.Port)
	if err != nil {
		panic("Error run")
	}
}
