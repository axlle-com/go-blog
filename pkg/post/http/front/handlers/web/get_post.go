package web

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (c *postController) GetPost(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if !isValidAlias(alias) {
		ctx.HTML(http.StatusNotFound, "404", gin.H{"title": "404 Not Found"})
		ctx.Abort()
		return
	}

	post, err := c.postService.GetByParam("alias", alias)
	if err != nil || post == nil {
		logger.Error(err)
		post = &models.Post{}
	}

	ctx.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title": "Home Page",
			"post":  post,
		},
	)
}

func isValidAlias(alias string) bool {
	validAliasPattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(validAliasPattern, alias)
	return matched
}
