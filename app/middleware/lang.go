package middleware

import (
	"strings"
	"time"

	i18nsvc "github.com/axlle-com/blog/app/service/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const CtxLocKey = "loc"

func NewLanguage(i18nService *i18nsvc.Service) *Language {
	return &Language{i18n: i18nService}
}

type Language struct {
	i18n *i18nsvc.Service
}

func (l *Language) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if l.i18n == nil {
			ctx.Next()
			return
		}

		accept, _ := l.i18n.NormalizeLang(ctx.GetHeader("Accept-Language"))

		qLangRaw := strings.TrimSpace(ctx.Query("lang"))
		qLang, qOk := l.i18n.NormalizeLang(qLangRaw)

		var cLang string
		if !qOk {
			if c, err := ctx.Request.Cookie("lang"); err == nil {
				if tag, ok := l.i18n.NormalizeLang(c.Value); ok {
					cLang = tag
				}
			}
		}

		loc := l.i18n.Localizer(qLang, cLang, accept)

		if qOk {
			cLang = qLang
			secure := ctx.Request.TLS != nil
			ctx.SetCookie(
				"lang",
				qLang,
				int((365 * 24 * time.Hour).Seconds()),
				"/",
				"",
				secure,
				false,
			)
		}

		if cLang == "" {
			cLang = accept
		}

		ctx.Set(CtxLocKey, loc)
		ctx.Set("lang", cLang)

		// функция T для шаблонов
		ctx.Set("T", buildT(loc))

		ctx.Next()
	}
}

// buildT создаёт функцию-переводчик для шаблонов
func buildT(loc *i18n.Localizer) func(id string, data map[string]any, n ...int) string {
	return func(id string, data map[string]any, n ...int) string {
		if loc == nil {
			return id
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
