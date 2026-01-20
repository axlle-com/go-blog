package models

import (
	"bytes"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/axlle-com/blog/app/middleware"
	"github.com/axlle-com/blog/app/models/contract"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-gonic/gin"
)

type BaseAjax struct{}

func (c *BaseAjax) GetID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0
	}

	return uint(id)
}

func (c *BaseAjax) GetUser(ctx *gin.Context) contract.User {
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

func (c *BaseAjax) GetAdmin(ctx *gin.Context) contract.User {
	userData, exists := ctx.Get("user")
	if !exists {
		return nil
	}

	u, ok := userData.(user.User)
	if !ok {
		return nil
	}

	return &u
}

// GetT достаёт функцию перевода из контекста
func (c *BaseAjax) GetT(ctx *gin.Context) func(id string, data map[string]any, n ...int) string {
	if v, ok := ctx.Get(middleware.CtxTKey); ok {
		if f, ok := v.(func(string, map[string]any, ...int) string); ok {
			return f
		}
	}

	// крайний случай — вернуть ключ
	return func(id string, _ map[string]any, _ ...int) string { return id }
}

// T — удобный шорткат для простых переводов (без плюрализации)
func (c *BaseAjax) T(ctx *gin.Context, id string, data ...map[string]any) string {
	t := c.GetT(ctx)
	var d map[string]any
	if len(data) > 0 {
		d = data[0]
	}
	return t(id, d)
}

// PrepareTemplateData приводит obj к единому виду gin.H и добавляет T (+ опционально settings.T)
func (c *BaseAjax) PrepareTemplateData(ctx *gin.Context, obj any) gin.H {
	tFunc := c.GetT(ctx)

	var data gin.H
	switch value := obj.(type) {
	case gin.H:
		data = value
	case map[string]any:
		data = gin.H(value)
	default:
		data = gin.H{"data": value}
	}

	// T в корень
	data["T"] = tFunc

	// T внутрь settings (если settings есть и это map)
	if settings, ok := data["settings"]; ok {
		switch s := settings.(type) {
		case gin.H:
			s["T"] = tFunc
		case map[string]any:
			s["T"] = tFunc
		}
	}

	return data
}

// RenderHTML рендерит шаблон, автоматически добавляя T в данные
func (c *BaseAjax) RenderHTML(ctx *gin.Context, code int, tplName string, obj any) {
	data := c.PrepareTemplateData(ctx, obj)
	ctx.HTML(code, tplName, data)
}

// RenderView рендерит HTML в строку (для ajax/partials), не трогая реальный response body
func (c *BaseAjax) RenderView(view string, data map[string]any, ctx *gin.Context) string {
	var buf bytes.Buffer

	originalWriter := ctx.Writer
	wrappedWriter := &ResponseWriterWrapper{
		ResponseWriter: originalWriter,
		Buffer:         &buf,
	}

	ctx.Writer = wrappedWriter
	defer func() { ctx.Writer = originalWriter }()

	c.RenderHTML(ctx, http.StatusOK, view, data)

	return c.removeWhitespaceBetweenTags(buf.String())
}

// removeWhitespaceBetweenTags — безопасная минификация только МЕЖДУ тегами
func (c *BaseAjax) removeWhitespaceBetweenTags(s string) string {
	re := regexp.MustCompile(`>\s+<`)
	s = re.ReplaceAllString(s, "><")
	return strings.TrimSpace(s)
}

func (c *BaseAjax) PaginatorFromQuery(ctx *gin.Context) contract.Paginator {
	return FromQuery(ctx.Request.URL.Query())
}
