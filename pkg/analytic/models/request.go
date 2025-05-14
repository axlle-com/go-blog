package models

import (
	errorsForm "github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func NewAnalyticRequest() *AnalyticRequest {
	return &AnalyticRequest{}
}

type AnalyticRequest struct {
	ID         uint   `json:"id"`
	From       string `gorm:"size:255null" json:"from" form:"from" binding:"omitempty"`
	To         string `gorm:"size:255;null" json:"to" form:"to" binding:"omitempty"`
	Subject    string `json:"subject" form:"subject" binding:"required,max=255"`
	Body       string `json:"body" form:"body" binding:"omitempty"`
	Attachment string `json:"attachment" form:"attachment" binding:"omitempty"`
	Viewed     string `json:"viewed"  form:"viewed" binding:"omitempty"`
}

func (m *AnalyticRequest) ValidateForm(ctx *gin.Context) (*AnalyticRequest, *errorsForm.Errors) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, &errorsForm.Errors{Message: "Форма не валидная!"}
	}
	if len(ctx.Request.PostForm) == 0 {
		return nil, &errorsForm.Errors{Message: "Форма не валидная!"}
	}
	if err := ctx.ShouldBind(&m); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}
	return m, nil
}

func (m *AnalyticRequest) ValidateJSON(ctx *gin.Context) (*AnalyticRequest, *errorsForm.Errors) {
	if err := ctx.ShouldBindJSON(&m); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}

	return m, nil
}
