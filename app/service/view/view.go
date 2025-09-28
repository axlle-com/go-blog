package view

import (
	"fmt"

	"github.com/axlle-com/blog/app/models/contracts"
)

type view struct {
	config contracts.Config
}

func NewView(config contracts.Config) contracts.View {
	return &view{config: config}
}

func (v *view) View(resource contracts.Resource) string {
	tpl := "index"
	if resource != nil {
		tpl = resource.GetTemplateName()
	}

	if v.config.Layout() == "" {
		return fmt.Sprintf("default.%s", tpl)
	}
	return fmt.Sprintf("%s.%s", v.config.Layout(), tpl)
}
