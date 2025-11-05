package view

import (
	"fmt"

	"github.com/axlle-com/blog/app/models/contract"
)

type view struct {
	config contract.Config
}

func NewView(config contract.Config) contract.View {
	return &view{config: config}
}

func (v *view) View(resource contract.Resource) string {
	tpl := "index"
	if resource != nil {
		tpl = resource.GetTemplateName()
	}

	if v.config.Layout() == "" {
		return fmt.Sprintf("default.%s", tpl)
	}
	return fmt.Sprintf("%s.%s", v.config.Layout(), tpl)
}

func (v *view) ViewStatic(name string) string {
	if name == "" {
		name = "index"
	}

	if v.config.Layout() == "" {
		return fmt.Sprintf("default.%s", name)
	}

	return fmt.Sprintf("%s.%s", v.config.Layout(), name)
}
