package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WebController interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	CreatePost(*gin.Context)
}

func NewWebController(r *gin.Engine) WebController {
	return &webController{engine: r}
}

type webController struct {
	engine *gin.Engine
}

func (c *webController) getID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0
	}
	return uint(id)
}

func (c *webController) getUser(ctx *gin.Context) *models.User {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return nil
	}
	user, ok := userData.(models.User)
	if !ok {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return nil
	}
	return &user
}
