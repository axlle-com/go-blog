package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	template "github.com/axlle-com/blog/pkg/template/repository"
	userRepo "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

func (c *controller) GetPosts(ctx *gin.Context) {
	start := time.Now()
	paginator := common.NewPaginator(ctx)

	user := c.getUser(ctx)
	if user == nil {
		return
	}

	users, err := userRepo.NewRepository().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	categories, err := repository.NewCategoryRepository().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := template.NewRepository().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	posts, total, err := repository.NewPostRepository().GetPaginate(paginator.GetPage(), paginator.GetPageSize())
	if err != nil {
		logger.New().Error(err)
	}
	log.Printf("Total time: %v", time.Since(start))
	paginator.SetTotal(total)
	ctx.HTML(
		http.StatusOK,
		"admin.posts",
		gin.H{
			"title":      "Страница постов",
			"posts":      posts,
			"categories": categories,
			"templates":  templates,
			"user":       user,
			"users":      users,
			"menu":       menu.NewMenu(ctx.FullPath()),
			"total":      total,
			"paginator":  paginator,
		},
	)
}

func GetPostsGo(ctx *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex

	userData, exists := ctx.Get("user")
	if !exists {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}
	user, ok := userData.(common.User)
	if !ok {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	paginator := common.NewPaginator(ctx)
	var users []common.User
	var categories []models.PostCategory
	var templates []common.Template
	var posts []models.PostResponse
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
		categories, err = repository.NewCategoryRepository().GetAll()
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
		posts, total, err = repository.NewPostRepository().GetPaginate(paginator.GetPage(), paginator.GetPageSize())
		if err != nil {
			mu.Lock()
			log.Println(err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	log.Printf("Total time: %v", time.Since(start))
	paginator.SetTotal(total)
	ctx.HTML(
		http.StatusOK,
		"admin.posts",
		gin.H{
			"title":      "Страница постов",
			"posts":      posts,
			"categories": categories,
			"templates":  templates,
			"user":       user,
			"users":      users,
			"menu":       menu.NewMenu(ctx.FullPath()),
			"total":      total,
			"paginator":  paginator,
		},
	)
}
