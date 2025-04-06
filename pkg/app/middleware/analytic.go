package middleware

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"strings"
	"time"
)

func Analytic() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger.Debug(detectDeviceType(c.GetHeader("User-Agent")))

		c.Next()

		logger.Debugf("Total time: %v", time.Since(start))
	}
}

func detectDeviceType(uaString string) string {
	ua := user_agent.New(uaString)

	switch {
	case ua.Bot():
		return "bot"
	case strings.Contains(uaString, "iPad"):
		return "tablet"
	case ua.Mobile():
		return "mobile"
	default:
		return "desktop"
	}
}
