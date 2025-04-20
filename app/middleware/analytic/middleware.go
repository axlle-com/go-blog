package analytic

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/analytic/provider"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"strconv"
	"strings"
	"time"
)

func NewAnalytic(
	queue contracts.Queue,
	analyticProvider provider.AnalyticProvider,
) *Analytic {
	return &Analytic{
		queue:            queue,
		analyticProvider: analyticProvider,
	}
}

type Analytic struct {
	queue            contracts.Queue
	analyticProvider provider.AnalyticProvider
}

func (a *Analytic) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ua := user_agent.New(c.GetHeader("User-Agent"))
		browserName, _ := ua.Browser()
		c.Set("device", detectDeviceType(c.GetHeader("User-Agent")))
		c.Set("browser", browserName)
		c.Set("os", ua.OS())
		c.Set("request_uuid", uuid.New().String())

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

		host := c.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}

		referer := c.GetHeader("Referer")
		if referer == "" {
			referer = c.GetHeader("Origin")
		}

		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			userUUID = c.GetString("guest_uuid")
		}

		evt := AnalyticsEvent{
			RequestUUID:      c.GetString("request_uuid"),
			UserUUID:         userUUID,
			Timestamp:        time.Now().UTC(),
			Method:           c.Request.Method,
			Host:             host,
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

		a.queue.Enqueue(NewAnalyticsJob(evt, a.analyticProvider), 0)
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
