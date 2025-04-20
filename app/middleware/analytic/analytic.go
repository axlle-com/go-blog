package analytic

import "time"

type AnalyticsEvent struct {
	RequestUUID      string        `json:"request_uuid"`
	UserUUID         string        `json:"user_uuid"`
	Timestamp        time.Time     `json:"timestamp"`
	Method           string        `json:"method"`
	Host             string        `json:"host"`
	Path             string        `json:"path"`
	Query            string        `json:"query,omitempty"`
	Status           int           `json:"status"`
	Latency          time.Duration `json:"latency"`
	IP               string        `json:"ip"`
	Country          string        `json:"country,omitempty"`
	City             string        `json:"city,omitempty"`
	Browser          string        `json:"browser,omitempty"`
	Device           string        `json:"device,omitempty"`
	OS               string        `json:"os,omitempty"`
	Language         string        `json:"lang,omitempty"`
	Referrer         string        `json:"referrer,omitempty"`
	ResolutionWidth  int           `json:"resolution_width,omitempty"`
	ResolutionHeight int           `json:"resolution_height,omitempty"`
	RequestSize      int64         `json:"req_size,omitempty"`
	ResponseSize     int64         `json:"resp_size,omitempty"`
	UTMCampaign      string        `json:"utm_campaign,omitempty"`
	UTMSource        string        `json:"utm_source,omitempty"`
	UTMMedium        string        `json:"utm_medium,omitempty"`
}
