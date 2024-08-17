package web

import (
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/menu"
	handlers "github.com/axlle-com/blog/pkg/post/http/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (c *webController) GetPosts(ctx *gin.Context) {
	start := time.Now()
	paginator := common.NewPaginator(ctx)

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	body := handlers.NewResponse(paginator).GetForWeb()
	body["title"] = "Страница постов"
	body["user"] = user
	body["menu"] = menu.NewMenu(ctx.FullPath())
	log.Printf("Total time: %v", time.Since(start))
	ctx.HTML(http.StatusOK, "admin.posts", body)
}
