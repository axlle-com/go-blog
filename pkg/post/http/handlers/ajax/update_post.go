package ajax

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) UpdatePost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	post, err := NewPostRepo().GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form, formError := NewPostRequest().ValidateForm(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors":  formError.Errors,
				"message": formError.Message,
			})
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	// TODO использовать post, заполнить из form
	form.ID = post.ID
	form.Image = post.Image
	form.UserID = post.UserID
	if err = form.UploadImageFile(ctx.Request); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	form.SetOriginal(post)
	err = NewPostRepo().Update(form)
	if err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Set("title", form.Title)
	galleries := gallery.Provider().SaveFromForm(ctx)
	for _, g := range galleries {
		err := g.Attach(form)
		if err != nil {
			logger.Error(err)
		}
	}

	categories, err := NewCategoryRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := template.Provider().GetAll()

	form.Galleries = gallery.Provider().GetAllForResource(form)

	data := gin.H{
		"categories": categories,
		"templates":  templates,
		"post":       form,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":       c.RenderView("admin.post_inner", data, ctx),
			"post":       form,
			"categories": categories,
			"templates":  templates,
		},
	})
}
