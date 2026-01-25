package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/service"
)

type MenuViewData struct {
	Menu       *models.Menu
	CurrentURL string
}

func NewMenuProvider(view contract.View, menuService *service.MenuService) app.MenuProvider {
	return &provider{
		view:        view,
		menuService: menuService,
	}
}

type provider struct {
	view        contract.View
	menuService *service.MenuService
}

func (p *provider) GetMenuString(id uint, url string) (string, error) {
	menu, err := p.menuService.GetMenuWithItems(id)
	if err != nil {
		return "", err
	}

	data := MenuViewData{
		Menu:       menu,
		CurrentURL: url,
	}

	return p.view.RenderToString(p.view.ViewResource(menu), data)
}
