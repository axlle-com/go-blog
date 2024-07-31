package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	UpdatePost(*gin.Context)
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	engine *gin.Engine
}

func (controller *controller) getID(c *gin.Context) uint {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		c.Abort()
	}
	return uint(id)
}

func (controller *controller) getUser(c *gin.Context) *models.User {
	userData, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return nil
	}
	user, ok := userData.(models.User)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return nil
	}
	return &user
}
