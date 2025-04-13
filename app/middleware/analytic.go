package middleware

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"strings"
	"time"
)

func Analytic() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger.Debug(detectDeviceType(c.GetHeader("User-Agent")))

		screenRes, err := c.Cookie("resolution")
		if err != nil {
			logger.Debugf("Cookie resolution not found %v", err)
		} else {
			parts := strings.Split(screenRes, ";")
			if len(parts) != 2 {
				logger.Warning("Invalid cookie format")
			} else {
				logger.Debugf("Resolution width: %v", parts[0])

				c.Set("resolution_width", parts[0])
				c.Set("resolution_height", parts[1])
			}
		}

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
