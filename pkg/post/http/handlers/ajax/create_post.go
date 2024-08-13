package ajax

import (
	"bytes"
	"fmt"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/service"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) createPost(ctx *gin.Context, ctr Container) {
	form := ctr.Request().ValidateForm(ctx)
	if form == nil {
		return
	}
	form.UserID = &c.getUser(ctx).ID
	err := form.UploadImageFile(ctx)
	if err := form.UploadImageFile(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		ctx.Abort()
		return
	}
	if err := ctr.Post().Create(form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			"url":  fmt.Sprintf("/admin/posts/%d", form.ID),
			"post": form,
		},
	})
}
