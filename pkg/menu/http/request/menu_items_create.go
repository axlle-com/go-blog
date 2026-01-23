package request

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewMenuItemsRequest() *MenuItemsRequest {
	return &MenuItemsRequest{}
}

type MenuItemsRequest struct {
	ID            *uint      `json:"id"`
	PublisherUUID *uuid.UUID `json:"publisher_uuid"`
	MenuID        uint       `json:"menu_id"`
	MenuItemID    *uint      `json:"menu_item_id,omitempty"`
	Path          string     `json:"path"`
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Ico           *string    `json:"ico,omitempty"`
	Sort          int        `json:"sort,omitempty"`
}

func (r *MenuItemsRequest) ValidateJSON(ctx *gin.Context) (*MenuItemsRequest, *errutil.Errors) {
	if err := ctx.ShouldBindJSON(&r); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return r, nil
}
