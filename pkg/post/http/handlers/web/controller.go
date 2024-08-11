package web

import (
	errorsForm "github.com/axlle-com/blog/pkg/common/errors"
	"github.com/axlle-com/blog/pkg/common/models"
	post "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
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
		// TODO
		ctx.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		ctx.Abort()
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

func (c *controller) ValidateForm(ctx *gin.Context) *post.Post {
	err := ctx.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Форма не валидная"})
		ctx.Abort()
		return nil
	}
	var form *post.Post
	if err := ctx.ShouldBind(&form); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  errBind,
			"message": "Ошибки валидации",
		})
		ctx.Abort()
		return nil
	}
	return form
}
