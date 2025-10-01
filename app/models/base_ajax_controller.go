package models

import (
	"bytes"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/app/models/contracts"
	user "github.com/axlle-com/blog/pkg/user/models"
)

type BaseAjax struct {
}

func (c *BaseAjax) GetID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0
	}
	return uint(id)
}

func (c *BaseAjax) GetUser(ctx *gin.Context) contracts.User {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return nil
	}
	u, ok := userData.(user.User)
	if !ok {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return nil
	}
	return &u
}

func (c *BaseAjax) RenderView(view string, data map[string]any, ctx *gin.Context) string {
	var buf bytes.Buffer
	originalWriter := ctx.Writer

	wrappedWriter := &ResponseWriterWrapper{
		ResponseWriter: ctx.Writer,
		Buffer:         &buf,
	}
	ctx.Writer = wrappedWriter
	ctx.HTML(http.StatusOK, view, data)

	ctx.Writer = originalWriter
	return c.removeWhitespaceBetweenTags(buf.String())
}

func (c *BaseAjax) removeWhitespaceBetweenTags(s string) string {
	re := regexp.MustCompile(`>\s+<`)
	compactHTML := re.ReplaceAllString(s, "><")
	return strings.TrimSpace(compactHTML)
}

func (c *BaseAjax) PaginatorFromQuery(ctx *gin.Context) contracts.Paginator {
	return PaginatorFromQuery(ctx.Request.URL.Query())
}
