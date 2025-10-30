package middleware

import (
	"net/http"
	"time"

	i18nsvc "github.com/axlle-com/blog/app/services/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const CtxLocKey = "loc"

func Language(i18n *i18nsvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		loc := i18n.Localizer(c.Request)

		// Если пришёл ?lang=xx — обновим cookie, чтобы закрепить язык
		if lang := c.Query("lang"); lang != "" {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "lang",
				Value:    lang,
				Path:     "/",
				MaxAge:   int((365 * 24 * time.Hour).Seconds()),
				HttpOnly: false,
				SameSite: http.SameSiteLaxMode,
			})
		}

		c.Set(CtxLocKey, loc)

		// Создаём функцию T для шаблонов
		tFunc := buildT(loc)
		c.Set("T", tFunc)

		c.Next()
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
