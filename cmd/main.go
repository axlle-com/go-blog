package main

import (
	"encoding/gob"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/middleware"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем или открываем файл для логов
	//f, err := os.OpenFile("/var/log/app/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//
	//// Настраиваем логгер
	//log.SetOutput(f)
	//
	//// Пример логов
	//log.Println("Это информационное сообщение")
	//log.Println("Это сообщение об ошибке")

	gob.Register(models.User{})
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
	router.Use(middleware.CurrentRouteMiddleware())

	//web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeApiRoutes(router)
	routes.InitializeWebRoutes(router)

	err = router.Run(cfg.Port)
	if err != nil {
		panic("Error run")
	}
}
