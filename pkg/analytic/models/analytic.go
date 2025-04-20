package models

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
	"time"
)

type Analytic struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	RequestUUID      *uuid.UUID    `gorm:"index" json:"request_uuid"`
	UserUUID         *uuid.UUID    `gorm:"index" json:"user_uuid"`
	Timestamp        time.Time     `gorm:"index" json:"timestamp"`
	Method           string        `json:"method"`
	Host             string        `json:"host"`
	Path             string        `json:"path"`
	Query            string        `json:"query,omitempty"`
	Status           int           `json:"status"`
	Latency          time.Duration `json:"latency"`
	IP               string        `json:"ip"`
	Country          *string       `json:"country,omitempty"`
	City             *string       `json:"city,omitempty"`
	Browser          string        `json:"browser,omitempty"`
	Device           string        `json:"device,omitempty"`
	OS               string        `json:"os,omitempty"`
	Language         string        `json:"lang,omitempty"`
	Referrer         string        `json:"referrer,omitempty"`
	ResolutionWidth  *int          `json:"resolution_width,omitempty"`
	ResolutionHeight *int          `json:"resolution_height,omitempty"`
	RequestSize      int64         `json:"req_size,omitempty"`
	ResponseSize     int64         `json:"resp_size,omitempty"`
	UTMCampaign      string        `json:"utm_campaign,omitempty"`
	UTMSource        string        `json:"utm_source,omitempty"`
	UTMMedium        string        `json:"utm_medium,omitempty"`

	User contracts.User `gorm:"-" json:"user"`
}

func (a *Analytic) GetTable() string {
	return "analytics"
}

func (a *Analytic) GetRequestUUID() string {
	if a.RequestUUID == nil {
		return ""
	}
	return a.RequestUUID.String()
}
