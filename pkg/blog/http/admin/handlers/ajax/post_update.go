package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/app/logger"
	. "github.com/axlle-com/blog/pkg/blog/http/admin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (c *postController) UpdatePost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	found, err := c.postService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

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

	form.ID = strconv.Itoa(int(id))
	form.UUID = found.UUID.String()
	post, err := c.postService.SaveFromRequest(form, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	tags, err := c.tagCollectionService.GetAll()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	templates := c.template.GetAll()
	infoBlocks := c.infoBlock.GetAll()

	data := gin.H{
		"tags":       tags,
		"categories": categories,
		"templates":  templates,
		"post":       post,
		"collection": gin.H{
			"infoBlocks":         infoBlocks,
			"ifoBlockCollection": post.InfoBlocks,
			"relationURL":        post.AdminURL(),
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"view": c.RenderView("admin.post_inner", data, ctx),
			"url":  fmt.Sprintf("/admin/posts/%d", post.ID),
			"post": post,
		},
	})
}
