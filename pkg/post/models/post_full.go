package models

type PostFull struct {
	Post
	UserFirstName      string  `json:"user_first_name"`
	UserLastName       string  `json:"user_last_name"`
	CategoryTitle      *string `json:"category_title"`
	CategoryTitleShort *string `json:"category_title_short"`
	TemplateTitle      *string `json:"template_title"`
	TemplateName       *string `json:"template_name"`
}

func (p *PostFull) Date() string {
	if p.CreatedAt == nil {
		return ""
	}
	return p.CreatedAt.Format("02.01.2006 15:04:05")
}

func (p *PostFull) GetTemplateTitle() string {
	if p.TemplateTitle == nil {
		return ""
	}
	return *p.TemplateTitle
}

func (p *PostFull) GetCategoryTitleShort() string {
	if p.CategoryTitleShort == nil {
		if p.CategoryTitle == nil {
			return ""
		}
		return *p.CategoryTitle
	}
	return *p.CategoryTitleShort
}
