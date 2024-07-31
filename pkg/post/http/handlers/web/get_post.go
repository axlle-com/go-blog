package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/menu"
	postRepo "github.com/axlle-com/blog/pkg/post/repository"
	postCategory "github.com/axlle-com/blog/pkg/post_category/repository"
	template "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (controller *controller) GetPost(c *gin.Context) {
	id := controller.getID(c)
	if id == 0 {
		return
	}

	user := controller.getUser(c)
	if user == nil {
		return
	}

	post, err := postRepo.NewRepository().GetByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin.404", gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		return
	}

	post.Galleries = gallery.NewProvider().GetAllForResource(post)

	categories, err := postCategory.NewRepository().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := template.NewRepository().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}
	log.Println(post.Galleries)
	c.HTML(
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":      "Страница поста",
			"user":       user,
			"categories": categories,
			"templates":  templates,
			"menu":       menu.NewMenu(c.FullPath()),
			"post":       post,
		},
	)
}
