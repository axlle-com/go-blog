package models

func NewTemplateFilter() *TemplateFilter {
	return &TemplateFilter{}
}

type TemplateFilter struct {
	ID           *uint   `json:"id" form:"id" binding:"omitempty"`
	UserID       *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Title        *string `json:"title" form:"title" binding:"omitempty"`
	Name         *string `json:"name" form:"name" binding:"omitempty"`
	Theme        *string `json:"theme" form:"theme" binding:"omitempty"`
	ResourceName *string `json:"resource_name" form:"resource_name" binding:"omitempty"`
	Date         *string `json:"date" form:"date" binding:"omitempty"`
}

func (f *TemplateFilter) SetResourceName(name string) *TemplateFilter {
	f.ResourceName = &name

	return f
}
