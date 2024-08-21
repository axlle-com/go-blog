package ajax

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	. "github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) filterPosts(ctx *gin.Context, ctr Container) {
	filter := ctr.Filter().ValidateForm(ctx)
	if filter == nil {
		return
	}

	paginator := ctr.Paginator(ctx)
	paginator.AddQueryString(string(filter.GetQueryString()))
	posts, err := ctr.Post().GetPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	categories, err := ctr.Category().GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := template.Provider().GetAll()
	users := user.Provider().GetAll()
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
	data["url"] = filter.GetURL()
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
