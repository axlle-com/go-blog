package web

import (
	"github.com/axlle-com/blog/app/logger"
	mApp "github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
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
	paginator := mApp.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	users := c.userProvider.GetAll()

	temp, err := c.messageCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Errorf("[GetMessages] Error: %v", err)
	}
	messages := c.messageCollectionService.Aggregates(temp)

	cnt, err := c.messageCollectionService.CountByField("viewed", false)
	if err != nil {
		logger.Errorf("[GetMessages] Error: %v", err)
	}

	ctx.HTML(http.StatusOK, "admin.messages", gin.H{
		"title":     "Страница сообщений",
		"message":   empty,
		"messages":  messages,
		"unviewed":  cnt,
		"users":     users,
		"paginator": paginator,
		"filter":    filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      models2.NewMenu(ctx.FullPath()),
		},
	})
}
