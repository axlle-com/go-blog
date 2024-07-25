package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	postRepo "github.com/axlle-com/blog/pkg/post/repository"
	postCategory "github.com/axlle-com/blog/pkg/post_category/repository"
	template "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetPostReserve(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin.404", gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		return
	}

	userData, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	user, ok := userData.(models.User)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	post, err := postRepo.NewRepository().GetByID(uint(id))
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "admin.404", gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		return
	}

	categories, err := postCategory.NewRepository().GetAll()
	if err != nil {
		log.Println(err)
	}

	templates, err := template.NewRepository().GetAllTemplates()
	if err != nil {
		log.Println(err)
	}

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