package analytic

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/analytic/provider"
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
	return func(ctx *gin.Context) {
		start := time.Now()

		ua := user_agent.New(ctx.GetHeader("User-Agent"))
		browserName, _ := ua.Browser()
		ctx.Set("device", detectDeviceType(ctx.GetHeader("User-Agent")))
		ctx.Set("browser", browserName)
		ctx.Set("os", ua.OS())
		ctx.Set("request_uuid", uuid.New().String())

		if res, err := ctx.Cookie("resolution"); err == nil {
			if p := strings.Split(res, ";"); len(p) == 2 {
				if w, _ := strconv.Atoi(p[0]); w > 0 {
					ctx.Set("resolution_width", w)
				}
				if h, _ := strconv.Atoi(p[1]); h > 0 {
					ctx.Set("resolution_height", h)
				}
			}
		}

		ctx.Next()

		host := ctx.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = ctx.Request.Host
		}

		referer := ctx.GetHeader("Referer")
		if referer == "" {
			referer = ctx.GetHeader("Origin")
		}

		userUUID := ctx.GetString("user_uuid")
		if userUUID == "" {
			userUUID = ctx.GetString("guest_uuid")
		}

		evt := AnalyticsEvent{
			RequestUUID:      ctx.GetString("request_uuid"),
			UserUUID:         userUUID,
			Timestamp:        time.Now().UTC(),
			Method:           ctx.Request.Method,
			Host:             host,
			Path:             ctx.FullPath(),
			Query:            ctx.Request.URL.RawQuery,
			Status:           ctx.Writer.Status(),
			Latency:          time.Since(start),
			IP:               ctx.ClientIP(),
			OS:               ctx.GetString("os"),
			Browser:          ctx.GetString("browser"),
			Device:           ctx.GetString("device"),
			Language:         ctx.GetString("lang"),
			Referrer:         referer,
			ResolutionWidth:  ctx.GetInt("resolution_width"),
			ResolutionHeight: ctx.GetInt("resolution_height"),
			RequestSize:      ctx.Request.ContentLength,
			ResponseSize:     int64(ctx.Writer.Size()),
			UTMCampaign:      ctx.Query("utm_campaign"),
			UTMSource:        ctx.Query("utm_source"),
			UTMMedium:        ctx.Query("utm_medium"),
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
