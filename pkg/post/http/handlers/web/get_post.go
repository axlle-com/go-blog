package web

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	post, err := c.postService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	post.Galleries = c.gallery.GetForResource(post)

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":        "Страница поста",
			"userProvider": user,
			"categories":   categories,
			"templates":    templates,
			"menu":         models.NewMenu(ctx.FullPath()),
			"post":         post,
		},
	)
}
