package ajax

import (
	"github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	updatePost(*gin.Context, Container)
	UpdatePostHandler() gin.HandlerFunc
	CreatePostHandler() gin.HandlerFunc
	DeletePostHandler() gin.HandlerFunc
	DeletePostImageHandler() gin.HandlerFunc
	DeletePost(*gin.Context, Container)
	DeletePostImage(*gin.Context, Container)
	createPost(*gin.Context, Container)
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	engine *gin.Engine
}

func (c *controller) getID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0
	}
	return uint(id)
}

func (c *controller) getUser(ctx *gin.Context) *models.User {
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
