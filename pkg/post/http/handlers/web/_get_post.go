package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	"github.com/axlle-com/blog/pkg/post/repository"
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
		c.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
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

	post, err := repository.NewPostRepository().GetByID(uint(id))
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	categories, err := repository.NewCategoryRepository().GetAll()
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
