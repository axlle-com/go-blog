package web

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

func (c *postController) FindByAlias(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if !isValidAlias(alias) {
		c.Render404(ctx, c.view.ViewStatic("404"), nil)
		return
	}

	post, _ := c.postService.FindByParam("alias", alias)
	if post != nil {
		c.RenderPost(ctx, post)
		return
	}

	c.Render404(ctx, c.view.ViewStatic("404"), nil)
	return
}

func isValidAlias(alias string) bool {
	validAliasPattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(validAliasPattern, alias)
	return matched
}
