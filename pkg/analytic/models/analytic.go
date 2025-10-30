package models

import (
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type Analytic struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	RequestUUID      *uuid.UUID `gorm:"index" json:"request_uuid"`
	UserUUID         *uuid.UUID `gorm:"index" json:"user_uuid"`
	Timestamp        time.Time  `gorm:"index" json:"timestamp"`
	Method           string     `gorm:"index" json:"method"`
	Host             string     `gorm:"index" json:"host"`
	Path             string     `gorm:"index" json:"path"`
	Query            string     `json:"query,omitempty"`
	Status           int        `gorm:"index" json:"status"`
	Latency          int64      `json:"latency"` // Milliseconds
	IP               string     `gorm:"index" json:"ip"`
	Country          *string    `gorm:"index" json:"country,omitempty"`
	City             *string    `gorm:"index" json:"city,omitempty"`
	Browser          string     `gorm:"index" json:"browser,omitempty"`
	Device           string     `gorm:"index" json:"device,omitempty"`
	OS               string     `gorm:"index" json:"os,omitempty"`
	Language         string     `gorm:"index" json:"lang,omitempty"`
	Referrer         string     `gorm:"index" json:"referrer,omitempty"`
	ResolutionWidth  *int       `gorm:"index" json:"resolution_width,omitempty"`
	ResolutionHeight *int       `json:"resolution_height,omitempty"`
	RequestSize      int64      `json:"req_size,omitempty"`  // KB
	ResponseSize     int64      `json:"resp_size,omitempty"` // KB
	UTMCampaign      string     `json:"utm_campaign,omitempty"`
	UTMSource        string     `json:"utm_source,omitempty"`
	UTMMedium        string     `json:"utm_medium,omitempty"`

	User contract.User `gorm:"-" json:"user"`
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
