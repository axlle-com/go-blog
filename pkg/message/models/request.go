package models

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func NewMessageRequest() *MessageRequest {
	return &MessageRequest{}
}

type MessageRequest struct {
	ID         uint   `json:"id"`
	From       string `gorm:"size:255null" json:"from" form:"from" binding:"omitempty"`
	To         string `gorm:"size:255;null" json:"to" form:"to" binding:"omitempty"`
	Subject    string `json:"subject" form:"subject" binding:"required,max=255"`
	Body       string `json:"body" form:"body" binding:"omitempty"`
	Attachment string `json:"attachment" form:"attachment" binding:"omitempty"`
	Viewed     string `json:"viewed"  form:"viewed" binding:"omitempty"`
}

func (m *MessageRequest) ValidateForm(ctx *gin.Context) (*MessageRequest, *errutil.Errors) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, &errutil.Errors{Message: "Форма не валидная!"}
	}

	if len(ctx.Request.PostForm) == 0 {
		return nil, &errutil.Errors{Message: "Форма не валидная!"}
	}

	if err := ctx.ShouldBind(&m); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return m, nil
}

func (m *MessageRequest) ValidateJSON(ctx *gin.Context) (*MessageRequest, *errutil.Errors) {
	if err := ctx.ShouldBindJSON(&m); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return m, nil
}
