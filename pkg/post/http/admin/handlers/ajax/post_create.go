package ajax

import (
	"fmt"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	. "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) CreatePost(ctx *gin.Context) {
	form, formError := NewPostRequest().ValidateJSON(ctx)
	if form == nil {
		if formError != nil {
			ctx.JSON(
				http.StatusBadRequest,
				response.Fail(http.StatusBadRequest, formError.Message, formError.Errors),
			)
			ctx.Abort()
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	post, err := c.postService.SaveFromRequest(form, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
		return
	}

	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
	infoBlocks := c.infoBlock.GetAll()

	data := response.Body{
		"categories": categories,
		"templates":  templates,
		"post":       post,
		"collection": gin.H{
			"infoBlocks":         infoBlocks,
			"ifoBlockCollection": post.InfoBlocks,
			"relationURL":        post.AdminURL(),
		},
	}
	ctx.JSON(
		http.StatusCreated,
		response.Created(
			response.Body{
				"view": c.RenderView("admin.post_inner", data, ctx),
				"url":  fmt.Sprintf("/admin/posts/%d", post.ID),
				"post": post,
			},
			"Запись создана",
		),
	)
}
