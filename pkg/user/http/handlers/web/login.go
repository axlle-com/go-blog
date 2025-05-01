package web

import (
	. "github.com/axlle-com/blog/app/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
)

func (c *controller) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID != nil {
		ctx.Redirect(http.StatusFound, "/")
		ctx.Abort()
		return
	}

	flashes := session.Flashes()
	errorMessages := ParseFlashes(flashes)
	err := session.Save()
	if err != nil {
		log.Println(err)
	}
	ctx.HTML(
		http.StatusOK,
		"admin.login",
		gin.H{
			"Title":  "Авторизация",
			"Errors": errorMessages,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
			},
		},
	)
}
