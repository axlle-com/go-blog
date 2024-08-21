package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/menu"
	. "github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (c *webController) getPosts(ctx *gin.Context, ctr Container) {
	start := time.Now()
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter := ctr.Filter().ValidateQuery(ctx)
	paginator := ctr.Paginator(ctx)
	paginator.AddQueryString(string(filter.GetQueryString()))
	templates := template.Provider().GetAll()
	users := userProvider.Provider().GetAll()
	categories, err := ctr.Category().GetAll()
	if err != nil {
		logger.Error(err)
	}

	posts, err := ctr.Post().GetPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	log.Printf("Total time: %v", time.Since(start))
	ctx.HTML(http.StatusOK, "admin.posts", gin.H{
		"title":      "Страница постов",
		"user":       user,
		"posts":      posts,
		"categories": categories,
		"templates":  templates,
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
		"menu":       menu.NewMenu(ctx.FullPath()),
	})
}
