package web

import (
	"bytes"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/post/repository"
	templateRepo "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) CreatePost(ctx *gin.Context) {
	form := c.ValidateForm(ctx)
	if form == nil {
		return
	}

	postRepo := repository.NewPostRepository()
	if err := postRepo.Create(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categories, err := repository.NewCategoryRepository().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := templateRepo.NewRepository().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	data := gin.H{
		"categories": categories,
		"templates":  templates,
		"post":       form,
	}

	var buf bytes.Buffer
	originalWriter := ctx.Writer

	wrappedWriter := &ResponseWriterWrapper{
		ResponseWriter: ctx.Writer,
		Buffer:         &buf,
	}
	ctx.Writer = wrappedWriter
	ctx.HTML(http.StatusOK, "admin.postInner", data)

	ctx.Writer = originalWriter
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view": buf.String(),
			"post": form,
		},
	})
}
