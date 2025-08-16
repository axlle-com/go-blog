package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/gin-gonic/gin"
)

func (c *messageController) GetMessages(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := models.NewMessageFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}

	empty := &models.Message{}
	paginator := app.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	temp, err := c.messageCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[GetMessages] Error: %v", err)
	}
	messages := c.messageCollectionService.Aggregates(temp)

	cnt, err := c.messageCollectionService.CountByField("viewed", false)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[GetMessages] Error: %v", err)
	}

	data := response.Body{
		"message":   empty,
		"messages":  messages,
		"paginator": paginator,
		"filter":    filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view":     c.RenderView("admin.messages_inner", data, ctx),
				"messages": messages,
				"unviewed": cnt,
			},
			"",
			nil,
		),
	)
}
