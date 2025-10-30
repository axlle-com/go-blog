package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Setting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Namespace string         `gorm:"index:ux_ns_key_scope,unique;not null" json:"namespace"`
	Key       string         `gorm:"index:ux_ns_key_scope,unique;not null" json:"key"`
	Type      string         `gorm:"not null" json:"type"`
	Value     datatypes.JSON `gorm:"type:jsonb;not null;default:'null'::jsonb" json:"value"`
	Scope     string         `gorm:"index:ux_ns_key_scope,unique;not null;default:global" json:"scope"`
	Sort      int            `gorm:"not null;default:100" json:"sort"`
	OwnerUUID *uuid.UUID     `json:"owner_uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Setting) TableName() string { return "settings" }
