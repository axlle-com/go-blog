package model

import (
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/google/uuid"
)

type Analytic struct {
	RequestUUID      string    `json:"request_uuid"`
	UserUUID         string    `json:"user_uuid"`
	Timestamp        time.Time `json:"timestamp"`
	Method           string    `json:"method"`
	Host             string    `json:"host"`
	Path             string    `json:"path"`
	Query            string    `json:"query,omitempty"`
	Status           int       `json:"status"`
	Latency          int64     `json:"latency"`
	IP               string    `json:"ip"`
	Country          *string   `json:"country,omitempty"`
	City             *string   `json:"city,omitempty"`
	Browser          string    `json:"browser,omitempty"`
	Device           string    `json:"device,omitempty"`
	OS               string    `json:"os,omitempty"`
	Language         string    `json:"lang,omitempty"`
	Referrer         string    `json:"referrer,omitempty"`
	ResolutionWidth  *int      `json:"resolution_width,omitempty"`
	ResolutionHeight *int      `json:"resolution_height,omitempty"`
	RequestSize      int64     `json:"req_size,omitempty"`
	ResponseSize     int64     `json:"resp_size,omitempty"`
	UTMCampaign      string    `json:"utm_campaign,omitempty"`
	UTMSource        string    `json:"utm_source,omitempty"`
	UTMMedium        string    `json:"utm_medium,omitempty"`
}

func (m *Analytic) Model() *models.Analytic {
	var (
		reqUUID  *uuid.UUID
		userUUID *uuid.UUID
	)

	if m.RequestUUID != "" {
		if u, err := uuid.Parse(m.RequestUUID); err != nil {
			logger.Errorf("[analytic][Model] invalid request_uuid: %v", err)
		} else {
			reqUUID = &u
		}
	}

	if m.UserUUID != "" {
		if u, err := uuid.Parse(m.UserUUID); err != nil {
			logger.Errorf("[analytic][Model] invalid user_uuid: %v", err)
		} else {
			userUUID = &u
		}
	}

	return &models.Analytic{
		RequestUUID:      reqUUID,
		UserUUID:         userUUID,
		Timestamp:        m.Timestamp,
		Method:           m.Method,
		Host:             m.Host,
		Path:             m.Path,
		Query:            m.Query,
		Status:           m.Status,
		Latency:          m.Latency,
		IP:               m.IP,
		Country:          m.Country,
		City:             m.City,
		Browser:          m.Browser,
		Device:           m.Device,
		OS:               m.OS,
		Language:         m.Language,
		Referrer:         m.Referrer,
		ResolutionWidth:  m.ResolutionWidth,
		ResolutionHeight: m.ResolutionHeight,
		RequestSize:      m.RequestSize,
		ResponseSize:     m.ResponseSize,
		UTMCampaign:      m.UTMCampaign,
		UTMSource:        m.UTMSource,
		UTMMedium:        m.UTMMedium,
	}
}
