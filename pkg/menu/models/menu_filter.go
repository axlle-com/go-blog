package models

func NewMenuFilter() *MenuFilter {
	return &MenuFilter{}
}

type MenuFilter struct {
	ID         *uint   `json:"id" form:"id" binding:"omitempty"`
	MenuID     *uint   `json:"menu_id" form:"menu_id" binding:"omitempty"`
	MenuItemID *uint   `json:"menu_item_id" form:"menu_item_id" binding:"omitempty"`
	Title      *string `json:"title" form:"title" binding:"omitempty"`

	array map[string]string // TODO map[string][]string
}

func (p *MenuFilter) SetMenuItemID(id uint) *MenuFilter {
	p.MenuItemID = &id
	return p
}

func (p *MenuFilter) SetMenuID(id uint) *MenuFilter {
	p.MenuID = &id
	return p
}

func (f *MenuFilter) SetMap(array map[string]string) {
	f.array = array
}

func (f *MenuFilter) GetMap() map[string]string {
	if f.array == nil {
		return nil
	}
	return f.array
}
