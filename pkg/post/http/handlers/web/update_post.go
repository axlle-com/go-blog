package web

import (
	"bytes"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/post/repository"
	templateRepo "github.com/axlle-com/blog/pkg/template/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) UpdatePost(ctx *gin.Context) {
	id := c.getID(ctx)
	if id == 0 {
		return
	}
	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	form := c.ValidateForm(ctx)
	if form == nil {
		return
	}
	form.ID = post.ID

	err = postRepo.Update(form)
	if err != nil {
		logger.New().Error(err)
	}
	ctx.Set("title", form.Title)
	galleries := service.SaveFromForm(ctx)
	for _, gallery := range galleries {
		err := gallery.Attach(form)
		if err != nil {
			logger.New().Error(err)
		}
	}

	categories, err := repository.NewCategoryRepository().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := templateRepo.NewRepository().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	form.Galleries = provider.NewProvider().GetAllForResource(form)

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

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	return rw.Buffer.Write(data)
}
