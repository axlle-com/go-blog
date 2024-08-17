package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func GetPostFront(c *gin.Context) {
	alias := c.Param("alias")
	if !isValidAlias(alias) {
		c.HTML(http.StatusNotFound, "404", gin.H{"title": "404 Not Found"})
		c.Abort()
		return
	}

	id := c.Param("id")
	h := db.GetDB()
	var post models.Post

	if result := h.First(&post, id); result.Error != nil {
		c.HTML(http.StatusNotFound, "404", gin.H{"title": "404 Not Found"})
		return
	}

	c.HTML(
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
