package web

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	mApp "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *messageController) DeleteMessage(ctx *gin.Context) {
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}

	message, err := c.messageService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"title": "404 Not Found"})
		return
	}

	err = c.messageService.Delete(message)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	empty := &models.Message{}
	paginator := mApp.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	temp, err := c.messageCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Errorf("[GetMessages] Error: %v", err)
	}
	messages := c.messageCollectionService.Aggregates(temp)

	data := response.Body{
		"messages":  messages,
		"paginator": paginator,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view":      c.RenderView("admin.messages_inner", data, ctx),
				"messages":  messages,
				"paginator": paginator,
			},
			"",
			nil,
		),
	)
}
