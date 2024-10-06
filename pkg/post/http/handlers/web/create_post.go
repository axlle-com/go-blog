package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/menu"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *webController) CreatePost(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	post := &Post{}
	categories, err := CategoryRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates, err := template.NewRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}
	ctx.HTML(
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":      "Страница поста",
			"user":       user,
			"categories": categories,
			"templates":  templates,
			"menu":       menu.NewMenu(ctx.FullPath()),
			"post":       post,
		},
	)
}
