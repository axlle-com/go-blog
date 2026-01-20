package web

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderByURL(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if !isValidAlias(alias) {
		c.RenderHTML(
			ctx,
			http.StatusNotFound,
			c.view.View("error"),
			gin.H{
				"title":    "Page not found",
				"error":    "404",
				"settings": c.settings(ctx, nil),
			},
		)
		ctx.Abort()
		return
	}

	url := "/" + alias
	post, _ := c.postService.FindByParam("url", url)
	if post != nil {
		c.RenderPost(ctx, post)
		return
	}

	category, _ := c.categoryService.FindByParam("url", url)
	if category != nil {
		c.RenderCategory(ctx, category)
		return
	}

	c.RenderHTML(
		ctx,
		http.StatusNotFound,
		c.view.View("error"),
		gin.H{
			"title":    "Page not found",
			"error":    "404",
			"settings": c.settings(ctx, nil),
		},
	)
	ctx.Abort()
	return
}

func isValidAlias(alias string) bool {
	validAliasPattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(validAliasPattern, alias)
	return matched
}
