package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.message.server_error")})
		return
	}

	empty := &models.Message{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	users := c.api.User.GetAll()

	temp, err := c.messageCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[GetMessages] Error: %v", err)
	}
	messages := c.messageCollectionService.Aggregates(temp)

	cnt, err := c.messageCollectionService.CountByField("viewed", false)
	if err != nil {
		logger.WithRequest(ctx).Errorf("[GetMessages] Error: %v", err)
	}

	c.RenderHTML(ctx, http.StatusOK, "admin.messages", gin.H{
		"title":     c.T(ctx, "ui.name.messages"),
		"message":   empty,
		"messages":  messages,
		"unviewed":  cnt,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
		},
	})
}
