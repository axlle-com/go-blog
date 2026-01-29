package models

import "html/template"

type CompanyInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Policy  string `json:"policy"`
}

func (ci *CompanyInfo) GetEmail() string {
	return ci.Email
}

func (ci *CompanyInfo) GetName() string {
	return ci.Name
}

func (ci *CompanyInfo) GetPhone() string {
	return ci.Phone
}

func (ci *CompanyInfo) GetAddress() string {
	return ci.Address
}

func (ci *CompanyInfo) GetPolicy() template.HTML {
	return template.HTML(ci.Policy)
}
