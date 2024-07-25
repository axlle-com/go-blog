package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	post "github.com/axlle-com/blog/pkg/post/models"
	postRepo "github.com/axlle-com/blog/pkg/post/repository"
	postCategory "github.com/axlle-com/blog/pkg/post_category/repository"
	template "github.com/axlle-com/blog/pkg/template/repository"
	userRepo "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

func (controller *controller) GetPosts(c *gin.Context) {
	start := time.Now()
	paginator := models.NewPaginator(c)

	user := controller.getUser(c)
	if user == nil {
		return
	}

	users, err := userRepo.NewRepository().GetAll()
	if err != nil {
		log.Println(err)
	}

	categories, err := postCategory.NewRepository().GetAll()
	if err != nil {
		log.Println(err)
	}

	templates, err := template.NewRepository().GetAllTemplates()
	if err != nil {
		log.Println(err)
	}

	posts, total, err := postRepo.NewRepository().GetPaginate(paginator.GetPage(), paginator.GetPageSize())
	if err != nil {
		log.Println(err)
	}
	log.Printf("Total time: %v", time.Since(start))
	paginator.SetTotal(total)
	c.HTML(
		http.StatusOK,
		"admin.posts",
		gin.H{
			"title":      "Страница постов",
			"posts":      posts,
			"categories": categories,
			"templates":  templates,
			"user":       user,
			"users":      users,
			"menu":       menu.NewMenu(c.FullPath()),
			"total":      total,
			"paginator":  paginator,
		},
	)
}

func GetPostsGo(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex

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

	paginator := models.NewPaginator(c)
	var users []models.User
	var categories []models.PostCategory
	var templates []models.Template
	var posts []post.Post
	var total int

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		users, err = userRepo.NewRepository().GetAll()
		if err != nil {
			mu.Lock()
			log.Println(err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		categories, err = postCategory.NewRepository().GetAll()
		if err != nil {
			mu.Lock()
			log.Println(err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		templates, err = template.NewRepository().GetAllTemplates()
		if err != nil {
			mu.Lock()
			log.Println(err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		posts, total, err = postRepo.NewRepository().GetPaginate(paginator.GetPage(), paginator.GetPageSize())
		if err != nil {
			mu.Lock()
			log.Println(err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	log.Printf("Total time: %v", time.Since(start))
	paginator.SetTotal(total)
	c.HTML(
		http.StatusOK,
		"admin.posts",
		gin.H{
			"title":      "Страница постов",
			"posts":      posts,
			"categories": categories,
			"templates":  templates,
			"user":       user,
			"users":      users,
			"menu":       menu.NewMenu(c.FullPath()),
			"total":      total,
			"paginator":  paginator,
		},
	)
}
