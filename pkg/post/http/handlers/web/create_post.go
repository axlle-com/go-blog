package web

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/menu"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) CreatePost(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	post := &models.Post{}
	categories, err := c.category.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
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
