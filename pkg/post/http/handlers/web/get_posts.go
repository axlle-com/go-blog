package web

import (
	post "github.com/axlle-com/blog/pkg/post/repository"
	postCategory "github.com/axlle-com/blog/pkg/post_category/repository"
	template "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetPosts(c *gin.Context) {
	posts, err := post.NewPostRepository().GetAllPosts()
	if err != nil {
		log.Println(err)
	}
	categories, err := postCategory.NewPostCategoryRepository().GetAllPostCategories()
	if err != nil {
		log.Println(err)
	}
	templates, err := template.NewTemplateRepository().GetAllTemplates()
	if err != nil {
		log.Println(err)
	}

	c.HTML(
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":      "Страница постов",
			"posts":      posts,
			"categories": categories,
			"templates":  templates,
		},
	)
}
