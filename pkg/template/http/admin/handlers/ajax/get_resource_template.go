package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateController) GetResourceTemplate(ctx *gin.Context) {
	template := ctx.Param("template")
	if template == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found"})
		ctx.Abort()
		return
	}

	template = models.NewResources().ResourceTemplate(template)
	if template == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found"})
		ctx.Abort()
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": template,
			},
			"",
			nil,
		),
	)
}
