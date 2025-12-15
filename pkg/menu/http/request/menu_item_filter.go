package request

import (
	"github.com/axlle-com/blog/app/errutil"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
)

func NewMenuItemFilter() *MenuItemFilter {
	return &MenuItemFilter{}
}

type MenuItemFilter struct {
	ID               *uint   `json:"id" form:"id" binding:"omitempty"`
	MenuID           *uint   `json:"menu_id" form:"menu_id" binding:"omitempty"`
	MenuItemID       *uint   `json:"menu_item_id" form:"menu_item_id" binding:"omitempty"`
	Query            *string `json:"query" form:"query" binding:"omitempty"`
	ForNotMenuItemID *uint   `json:"for_not_menu_item_id" form:"for_not_menu_item_id" binding:"omitempty"`
	app.Filter
}

func (p *MenuItemFilter) ValidateQuery(ctx *gin.Context) (*MenuItemFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *MenuItemFilter) ToFilter() *models.MenuItemFilter {
	filter := models.NewMenuItemFilter()
	filter.MenuItemID = p.MenuItemID
	filter.MenuID = p.MenuID
	filter.Title = p.Query
	filter.ForNotMenuItemID = p.ForNotMenuItemID
	filter.SetMap(p.GetMap())

	return filter
}

func (p *MenuItemFilter) SetMenuItemID(id uint) *MenuItemFilter {
	p.MenuItemID = &id
	return p
}

func (p *MenuItemFilter) SetMenuID(id uint) *MenuItemFilter {
	p.MenuID = &id
	return p
}
