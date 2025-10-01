package request

import (
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
)

func NeTemplateFilter() *TemplateFilter {
	return &TemplateFilter{}
}

type TemplateFilter struct {
	ID           *uint   `json:"id" form:"id" binding:"omitempty"`
	UserID       *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Title        *string `json:"title" form:"title" binding:"omitempty"`
	Name         *string `json:"name" form:"name" binding:"omitempty"`
	ResourceName *string `json:"resource_name" form:"resource_name" binding:"omitempty"`
	Date         *string `json:"date" form:"date" binding:"omitempty"`
	app.Filter
}

func (p *TemplateFilter) ValidateForm(ctx *gin.Context) (*TemplateFilter, *errutil.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *TemplateFilter) ValidateQuery(ctx *gin.Context) (*TemplateFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *TemplateFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *TemplateFilter) GetURL() string {
	return string("templates?" + p.GetQueryString())
}

func (p *TemplateFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}

func (p *TemplateFilter) ToFilter() *models.TemplateFilter {
	filter := models.NewTemplateFilter()
	filter.ID = p.ID
	filter.ResourceName = p.ResourceName

	return filter
}
