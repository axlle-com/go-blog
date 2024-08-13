package ajax

import (
	"bytes"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/service"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) updatePost(ctx *gin.Context, ctr Container) {
	id := c.getID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	post, err := ctr.Post().GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form := ctr.Request().ValidateForm(ctx)
	if form == nil {
		return
	}
	form.ID = post.ID
	form.Image = post.Image
	form.UserID = post.UserID
	if err = form.UploadImageFile(ctx); err != nil {
		logger.New().Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	form.SetOriginal(post)
	err = ctr.Post().Update(form)
	if err != nil {
		logger.New().Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Set("title", form.Title)
	galleries := service.SaveFromForm(ctx)
	for _, gallery := range galleries {
		err := gallery.Attach(form)
		if err != nil {
			logger.New().Error(err)
		}
	}

	categories, err := ctr.Category().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := ctr.Template().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	form.Galleries = ctr.GalleryProvider().GetAllForResource(form)

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
	ctx.HTML(http.StatusOK, "admin.post_inner", data)

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
