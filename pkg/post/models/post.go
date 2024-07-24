package models

import "github.com/axlle-com/blog/pkg/common/models"

type Post struct {
	models.Post
	UserFirstName      string  `json:"user_first_name"`
	UserLastName       string  `json:"user_last_name"`
	CategoryTitle      *string `json:"category_title"`
	CategoryTitleShort *string `json:"category_title_short"`
	TemplateTitle      *string `json:"template_title"`
	TemplateName       *string `json:"template_name"`
}

func (p *Post) Date() string {
	return p.CreatedAt.Format("02.01.2006 15:04:05")
}
