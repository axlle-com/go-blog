package main

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Error read config")
	}

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	router := gin.Default()

	h := db.Init(dbUrl)

	web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeWebRoutes(router, h)
	routes.InitializeApiRoutes(router, h)

	err = router.Run(port)
	if err != nil {
		panic("Error run")
	}
}
