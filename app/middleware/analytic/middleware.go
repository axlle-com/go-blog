package analytic

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"strconv"
	"strings"
	"time"
)

func NewAnalytic(queue contracts.Queue) *Analytic {
	return &Analytic{
		queue: queue,
	}
}

type Analytic struct {
	queue contracts.Queue
}

func (a *Analytic) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ua := user_agent.New(c.GetHeader("User-Agent"))
		browserName, _ := ua.Browser()
		c.Set("device", detectDeviceType(c.GetHeader("User-Agent")))
		c.Set("browser", browserName)
		c.Set("os", ua.OS())
		c.Set("user_uuid", "")

		if res, err := c.Cookie("resolution"); err == nil {
			if p := strings.Split(res, ";"); len(p) == 2 {
				if w, _ := strconv.Atoi(p[0]); w > 0 {
					c.Set("resolution_width", w)
				}
				if h, _ := strconv.Atoi(p[1]); h > 0 {
					c.Set("resolution_height", h)
				}
			}
		}

		c.Next()

		referer := c.GetHeader("Referer")
		if referer == "" {
			referer = c.GetHeader("Origin")
		}

		evt := AnalyticsEvent{
			RequestID:        c.GetString("request_id"),
			UserUUID:         c.GetString("user_uuid"),
			Timestamp:        time.Now().UTC(),
			Method:           c.Request.Method,
			Path:             c.FullPath(),
			Query:            c.Request.URL.RawQuery,
			Status:           c.Writer.Status(),
			Latency:          time.Since(start),
			IP:               c.ClientIP(),
			OS:               c.GetString("os"),
			Browser:          c.GetString("browser"),
			Device:           c.GetString("device"),
			Language:         c.GetString("lang"),
			Referrer:         referer,
			ResolutionWidth:  c.GetInt("resolution_width"),
			ResolutionHeight: c.GetInt("resolution_height"),
			RequestSize:      c.Request.ContentLength,
			ResponseSize:     int64(c.Writer.Size()),
			UTMCampaign:      c.Query("utm_campaign"),
			UTMSource:        c.Query("utm_source"),
			UTMMedium:        c.Query("utm_medium"),
		}

		a.queue.Enqueue(NewAnalyticsJob(evt), 0)

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
