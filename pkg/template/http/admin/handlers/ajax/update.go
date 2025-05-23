package ajax

import (
	mApp "github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateController) UpdateTemplate(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.templateService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form, formError := NewTemplateRequest().ValidateJSON(ctx)
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

	template, err := c.templateService.SaveFromRequest(form, found, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	data := gin.H{
		"templateModel": template,
		"resources":     mApp.NewResources().Resources(),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view":     c.RenderView("admin.template_inner", data, ctx),
			"url":      template.AdminURL(),
			"template": template,
		},
	})
}
