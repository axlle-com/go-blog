package ajax

import (
	"fmt"
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/http/admin/request"
	"github.com/gin-gonic/gin"
)

func (c *postController) CreatePost(ctx *gin.Context) {
	form, formError := request.NewPostRequest().ValidateJSON(ctx)
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

	post, err := c.postService.SaveFromRequest(form, nil, c.GetUser(ctx))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.Fail(http.StatusInternalServerError, err.Error(), nil),
		)
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

	infoBlocks := c.api.InfoBlock.GetAll()

	data := response.Body{
		"tags":       tags,
		"categories": categories,
		"templates":  c.templates(ctx),
		"post":       post,
		"collection": gin.H{
			"infoBlocks":          infoBlocks,
			"infoBlockCollection": post.InfoBlocks,
			"relationURL":         post.AdminURL(),
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
			c.T(ctx, "ui.success.record_created"),
		),
	)
}
