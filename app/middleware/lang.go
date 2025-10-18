package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Language() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang := ctx.GetHeader("Accept-Language")

		// оставим только первый язык, если их несколько
		if len(lang) > 0 {
			if i := strings.Index(lang, ","); i != -1 {
				lang = lang[:i]
			}
		}

		ctx.Set("lang", lang)
		ctx.Next()
	}
}
