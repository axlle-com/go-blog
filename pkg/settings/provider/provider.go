package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
	appProvider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/settings/service"
)

func NewProvider(
	companyInfoService *service.CompanyInfoService,
) appProvider.CompanyInfoProvider {
	return &provider{
		companyInfoService: companyInfoService,
	}
}

type provider struct {
	companyInfoService *service.CompanyInfoService
}

func (p *provider) GetCompanyInfo(ns, scope string) (contract.CompanyInfo, bool) {
	return p.companyInfoService.GetCompanyInfo(ns, scope)
}
