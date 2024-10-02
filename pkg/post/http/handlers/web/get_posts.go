package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/provider"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (c *webController) GetPosts(ctx *gin.Context) {
	start := time.Now()
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := NewPostFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}
	paginator := models.Paginator(ctx.Request.URL.Query())
	paginator.AddQueryString(string(filter.GetQueryString()))
	templates := template.Provider().GetAll()
	users := userProvider.Provider().GetAll()
	categories, err := NewCategoryRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}

	posts, err := NewPostRepo().GetPaginate(paginator, filter)
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
