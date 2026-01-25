package request

import (
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
)

func NewInfoBlockRequest() *InfoBlockRequest {
	return &InfoBlockRequest{}
}

type InfoBlockRequest struct {
	ID             *uint   `json:"id" form:"id" binding:"omitempty"`
	TemplateName   *string `json:"template_name" form:"template_name" binding:"omitempty"`
	UserID         *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	PostCategoryID *uint   `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string `json:"title" form:"title" binding:"omitempty"`
	Date           *string `json:"date" form:"date" binding:"omitempty"`
	app.Filter
}

func (p *InfoBlockRequest) ValidateQuery(ctx *gin.Context) (*InfoBlockRequest, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *InfoBlockRequest) PrintTemplateName() string {
	if p.TemplateName == nil {
		return ""
	}

	return *p.TemplateName
}

func (p *InfoBlockRequest) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *InfoBlockRequest) PrintPostCategoryID() uint {
	if p.PostCategoryID == nil {
		return 0
	}
	return *p.PostCategoryID
}

func (p *InfoBlockRequest) GetURL() string {
	return string("info-blocks?" + p.GetQueryString())
}

func (p *InfoBlockRequest) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}

func (p *InfoBlockRequest) ToFilter() *models.InfoBlockFilter {
	filter := models.NewInfoBlockFilter()
	filter.ID = p.ID
	filter.TemplateName = p.TemplateName
	filter.UserID = p.UserID
	filter.PostCategoryID = p.PostCategoryID
	filter.Title = p.Title
	filter.Date = p.Date

	return filter
}
