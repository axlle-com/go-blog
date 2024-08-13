package ajax

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/menu"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPost(ctx *gin.Context) {
	id := c.getID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	user := c.getUser(ctx)
	if user == nil {
		return
	}

	post, err := NewPostRepo().GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	post.Galleries = gallery.NewProvider().GetAllForResource(post)

	categories, err := NewCategoryRepo().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := template.NewRepo().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
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
