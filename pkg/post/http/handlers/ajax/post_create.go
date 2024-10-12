package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/logger"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) CreatePost(ctx *gin.Context) {
	form, formError := NewPostRequest().ValidateJSON(ctx)
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
	post, err := service.PostCreate(form, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusMultipleChoices, gin.H{"message": "err.Error()"})
	return

	//ctx.Set("title", form.Title)
	//galleries := gallery.Provider().SaveFromForm(ctx)
	//for _, g := range galleries {
	//	err := g.Attach(form)
	//	if err != nil {
	//		logger.Error(err)
	//	} else {
	//		form.Galleries = append(form.Galleries, g)
	//	}
	//}

	categories, err := CategoryRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := template.Provider().GetAll()

	data := gin.H{
		"categories": categories,
		"templates":  templates,
		"post":       post,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view": c.RenderView("admin.post_inner", data, ctx),
			"url":  fmt.Sprintf("/admin/posts/%d", post.ID),
			"post": post,
		},
	})
}
