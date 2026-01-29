package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	SettingTypeString = "string"
	SettingTypeBool   = "bool"
	SettingTypeJSON   = "json"

	CompanyEmailKey   = "email"
	CompanyNameKey    = "name"
	CompanyPhoneKey   = "phone"
	CompanyAddressKey = "address"
	PolicyKey         = "policy"
)

type Setting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Namespace string         `gorm:"not null" json:"namespace"`
	Key       string         `gorm:"not null" json:"key"`
	Type      string         `gorm:"not null" json:"type"`
	Value     datatypes.JSON `gorm:"type:jsonb;not null;default:'null'::jsonb" json:"value"`
	Scope     string         `gorm:"not null;default:global" json:"scope"`
	Sort      int            `gorm:"not null;default:100" json:"sort"`
	OwnerUUID *uuid.UUID     `json:"owner_uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Setting) GetTable() string {
	return "settings"
}
