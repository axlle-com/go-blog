package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *postController) GetPost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	post, err := c.postService.FindAggregateByID(id)
	if err != nil {
		logger.WithRequest(ctx).Error(err.Error())
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	tags, err := c.tagCollectionService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.post",
		gin.H{
			"title":      "Страница поста",
			"tags":       tags,
			"categories": categories,
			"templates":  c.templates(ctx),
			"post":       post,
			"collection": gin.H{
				"infoBlocks":          c.api.InfoBlock.GetAll(),
				"infoBlockCollection": post.InfoBlocks,
				"relationURL":         post.AdminURL(),
			},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
