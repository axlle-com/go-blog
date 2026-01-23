package provider

import "github.com/axlle-com/blog/app/models/contract"

type CompanyInfoProvider interface {
	GetCompanyInfo(ns, scope string) (contract.CompanyInfo, bool)
}
