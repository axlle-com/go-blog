package request

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func NewMenuRequest() *MenuRequest {
	return &MenuRequest{}
}

type MenuRequest struct {
	ID          uint    `json:"id"`
	TemplateID  *uint   `json:"template_id"`
	Title       string  `json:"title"`
	IsPublished *bool   `json:"is_published,omitempty"`
	IsMain      *bool   `json:"is_main,omitempty"`
	Ico         *string `json:"ico,omitempty"`
	Sort        int     `json:"sort,omitempty"`

	MenuItems []*MenuItemsRequest `json:"menu_items" form:"menu_items"`
}

func (r *MenuRequest) ValidateForm(ctx *gin.Context) (*MenuRequest, *errutil.Errors) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, &errutil.Errors{Message: "Форма не валидная!"}
	}

	if len(ctx.Request.PostForm) == 0 {
		return nil, &errutil.Errors{Message: "Форма не валидная!"}
	}

	if err := ctx.ShouldBind(&r); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return r, nil
}

func (r *MenuRequest) ValidateJSON(ctx *gin.Context) (*MenuRequest, *errutil.Errors) {
	if err := ctx.ShouldBindJSON(&r); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return r, nil
}
