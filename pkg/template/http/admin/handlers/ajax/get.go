package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateController) GetTemplate(ctx *gin.Context) {
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

	ctx.JSON(
		http.StatusCreated,
		response.OK(
			response.Body{
				"view":     nil,
				"url":      nil,
				"template": found,
			},
			"Запись создана",
			nil,
		),
	)
}
