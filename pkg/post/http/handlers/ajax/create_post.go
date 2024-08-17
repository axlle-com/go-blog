package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/logger"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) createPost(ctx *gin.Context, ctr Container) {
	form := ctr.Request().ValidateForm(ctx)
	if form == nil {
		return
	}
	form.UserID = &c.GetUser(ctx).ID
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
	galleries := ctr.Gallery().SaveFromForm(ctx)
	for _, gallery := range galleries {
		err := gallery.Attach(form)
		if err != nil {
			logger.Error(err)
		}
	}

	categories, err := ctr.Category().GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates, err := ctr.Template().GetAllTemplates()
	if err != nil {
		logger.Error(err)
	}

	data := gin.H{
		"categories": categories,
		"templates":  templates,
		"post":       form,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view": c.RenderView("admin.post_inner", data, ctx),
			"url":  fmt.Sprintf("/admin/posts/%d", form.ID),
			"post": form,
		},
	})
}
