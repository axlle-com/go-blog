package models

func NewMenuItemFilter() *MenuItemFilter {
	return &MenuItemFilter{}
}

type MenuItemFilter struct {
	ID               *uint   `json:"id" form:"id" binding:"omitempty"`
	MenuID           *uint   `json:"menu_id" form:"menu_id" binding:"omitempty"`
	MenuItemID       *uint   `json:"menu_item_id" form:"menu_item_id" binding:"omitempty"`
	Title            *string `json:"title" form:"title" binding:"omitempty"`
	ForNotMenuItemID *uint   `json:"for_not_menu_item_id" form:"for_not_menu_item_id" binding:"omitempty"`
	IDs              []uint  `json:"ids" form:"ids" binding:"omitempty"`

	array map[string]string // TODO map[string][]string
}

func (f *MenuItemFilter) SetMenuItemID(id uint) *MenuItemFilter {
	f.MenuItemID = &id
	return f
}

func (f *MenuItemFilter) SetMenuID(id uint) *MenuItemFilter {
	f.MenuID = &id
	return f
}

func (f *MenuItemFilter) SetMap(array map[string]string) {
	f.array = array
}

func (f *MenuItemFilter) GetMap() map[string]string {
	if f.array == nil {
		return nil
	}
	return f.array
}
