package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/gin-gonic/gin"
)

func (c *messageController) GetMessage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"title": "404 Not Found"})
		return
	}

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.message.server_error")})
		return
	}

	message, err := c.messageService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"title": "404 Not Found"})
		return
	}

	message.Viewed = true
	message, err = c.messageService.Update(message)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	message = c.messageService.Aggregate(message)

	empty := &models.Message{}
	paginator := c.PaginatorFromQuery(ctx)
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
		"message":   message,
		"paginator": paginator,
	}

	dataList := response.Body{
		"messages":  messages,
		"paginator": paginator,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view":     c.RenderView("admin.message_inner", data, ctx),
				"list":     c.RenderView("admin.messages_inner", dataList, ctx),
				"message":  message,
				"messages": messages,
				"unviewed": cnt,
			},
			"",
			nil,
		),
	)
}
