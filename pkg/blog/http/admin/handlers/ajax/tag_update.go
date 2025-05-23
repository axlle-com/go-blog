package ajax

import (
	. "github.com/axlle-com/blog/pkg/blog/http/admin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (c *tagController) Update(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.tagService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	form, formError := NewTagRequest().ValidateJSON(ctx)
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

	form.ID = strconv.Itoa(int(id))
	form.UUID = found.UUID.String()
	tag, err := c.tagService.SaveFromRequest(form, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	templates := c.template.GetAll()
	infoBlocks := c.infoBlock.GetAll()

	data := gin.H{
		"templates": templates,
		"tag":       tag,
		"collection": gin.H{
			"infoBlocks":         infoBlocks,
			"ifoBlockCollection": tag.InfoBlocks,
			"relationURL":        tag.AdminURL(),
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view": c.RenderView("admin.tag_inner", data, ctx),
			"url":  tag.AdminURL(),
			"post": tag,
		},
	})
}
