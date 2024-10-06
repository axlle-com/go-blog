package ajax

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePost(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	if err := PostRepo().Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	filter, validError := NewPostFilter().ValidateQuery(ctx)
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
	paginator := models.Paginator(ctx.Request.URL.Query())
	paginator.AddQueryString(string(filter.GetQueryString()))
	users := user.Provider().GetAll()
	templates := template.Provider().GetAll()
	categories, err := CategoryRepo().GetAll()
	if err != nil {
		logger.Error(err)
	}
	posts, err := PostRepo().GetPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}

	data := gin.H{
		"title":      "Страница постов",
		"posts":      posts,
		"categories": categories,
		"templates":  templates,
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
	}

	data["view"] = c.RenderView("admin.posts_inner", data, ctx)
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
