package models

import (
	"github.com/google/uuid"
)

func NewGalleryFilter() *GalleryFilter {
	return &GalleryFilter{}
}

type GalleryFilter struct {
	ID           *uint      `json:"id" form:"id" binding:"omitempty"`
	IDs          []uint     `json:"ids" form:"ids" binding:"omitempty"`
	TemplateID   *uint      `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID       *uint      `json:"user_id" form:"user_id" binding:"omitempty"`
	Title        *string    `json:"title" form:"title" binding:"omitempty"`
	ResourceUUID *uuid.UUID `json:"resource_uuid" form:"resource_uuid" binding:"omitempty"`
}
