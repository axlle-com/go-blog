package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *postController) CreatePost(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	post := &models.Post{}
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	tags, err := c.tagCollectionService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	templates := c.template.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":      "Страница поста",
			"categories": categories,
			"tags":       tags,
			"templates":  templates,
			"post":       post,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menu.NewMenu(ctx.FullPath()),
			},
		},
	)
}
