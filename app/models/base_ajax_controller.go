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
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (c *BaseAjax) RenderView(view string, data map[string]any, ctx *gin.Context) string {
	var buf bytes.Buffer
	originalWriter := ctx.Writer

	wrappedWriter := &ResponseWriterWrapper{
		ResponseWriter: ctx.Writer,
		Buffer:         &buf,
	}
	ctx.Writer = wrappedWriter
	c.RenderHTML(ctx, http.StatusOK, view, data)

	ctx.Writer = originalWriter
	return c.removeWhitespaceBetweenTags(buf.String())
}

// RenderHTML автоматически добавляет функцию перевода T в данные шаблона
func (c *BaseAjax) RenderHTML(ctx *gin.Context, code int, name string, obj any) {
	data := c.prepareTemplateData(ctx, obj)
	ctx.HTML(code, name, data)
}

// prepareTemplateData добавляет T в данные шаблона, если obj - это gin.H или map
func (c *BaseAjax) prepareTemplateData(ctx *gin.Context, obj any) any {
	return PrepareTemplateData(ctx, obj, c.BuildT(ctx))
}

func (c *BaseAjax) removeWhitespaceBetweenTags(s string) string {
	re := regexp.MustCompile(`>\s+<`)
	compactHTML := re.ReplaceAllString(s, "><")
	return strings.TrimSpace(compactHTML)
}

func (c *BaseAjax) PaginatorFromQuery(ctx *gin.Context) contract.Paginator {
	return PaginatorFromQuery(ctx.Request.URL.Query())
}

// getLoc получает локализатор из контекста.
func getLoc(ctx *gin.Context) *i18n.Localizer {
	v, ok := ctx.Get(middleware.CtxLocKey)
	if !ok {
		return nil
	}
	loc, _ := v.(*i18n.Localizer)
	return loc
}

// BuildT возвращает замыкание-переводчик, привязанное к текущему запросу.
func (c *BaseAjax) BuildT(ctx *gin.Context) func(id string, data map[string]any, n ...int) string {
	loc := getLoc(ctx)
	return func(id string, data map[string]any, n ...int) string {
		if loc == nil {
			return id // на крайний случай
		}
		cfg := &i18n.LocalizeConfig{MessageID: id, TemplateData: data}
		if len(n) > 0 {
			cfg.PluralCount = n[0]
		}
		s, err := loc.Localize(cfg)
		if err != nil || s == "" {
			return id
		}
		return s
	}
}

// PrepareTemplateData - универсальная функция для подготовки данных шаблона с автоматическим добавлением T
// Можно использовать в любых контроллерах, даже если они не наследуются от BaseAjax
func PrepareTemplateData(ctx *gin.Context, obj any, fallbackT func(id string, data map[string]any, n ...int) string) any {
	// Получаем T из контекста
	var tFunc func(id string, data map[string]any, n ...int) string
	if t, ok := ctx.Get("T"); ok {
		if f, ok := t.(func(id string, data map[string]any, n ...int) string); ok {
			tFunc = f
		}
	}

	// Если T не найдена, используем fallback
	if tFunc == nil {
		tFunc = fallbackT
	}

	// Добавляем T в данные и во вложенные settings
	switch v := obj.(type) {
	case gin.H:
		v["T"] = tFunc
		// Добавляем T и Title во вложенные settings, если они есть
		addTToSettings(v, tFunc)
		return v
	case map[string]any:
		v["T"] = tFunc
		// Добавляем T и Title во вложенные settings, если они есть
		addTToSettings(v, tFunc)
		return v
	default:
		// Если obj не map, создаём новый gin.H
		return gin.H{
			"T":    tFunc,
			"data": obj,
		}
	}
}

// addTToSettings добавляет T и Title во вложенные объекты settings
func addTToSettings(data map[string]any, tFunc func(id string, data map[string]any, n ...int) string) {
	if settings, ok := data["settings"]; ok {
		switch s := settings.(type) {
		case gin.H:
			s["T"] = tFunc
			// Добавляем Title из корневого контекста, если он есть
			if title, ok := data["title"]; ok {
				s["Title"] = title
			}
		case map[string]any:
			s["T"] = tFunc
			// Добавляем Title из корневого контекста, если он есть
			if title, ok := data["title"]; ok {
				s["Title"] = title
			}
		}
	}
}
