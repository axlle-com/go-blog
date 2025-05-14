package web

import (
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
		ctx.HTML(http.StatusNotFound, "404", gin.H{"title": "404 Not Found"})
		ctx.Abort()
		return
	}

	ctx.HTML(
		http.StatusOK,
		c.view.View(post),
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
