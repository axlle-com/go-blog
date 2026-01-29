package contract

import "html/template"

type CompanyInfo interface {
	GetEmail() string
	GetName() string
	GetPhone() string
	GetAddress() string
	GetPolicy() template.HTML
}
