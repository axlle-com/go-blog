package models

import (
	"html/template"
	"strings"
)

type AdminMenu struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"` // Ключ перевода или переведенный текст
	Ico      template.HTML `json:"ico"`
	IsActive bool          `json:"is_active"`
}

// NewMenu создает меню с ключами переводов
// Если передана функция перевода tFunc, то сразу переводит названия
// Если нет - возвращает ключи переводов, которые можно перевести в шаблоне
func NewMenu(currentRoute string, tFunc func(id string, data map[string]any, n ...int) string) []AdminMenu {
	var routes = []AdminMenu{
		{Path: "/admin/", Name: "ui.name.dashboard", Ico: template.HTML("<i data-feather=\"globe\"></i>")},
		{Path: "/admin/menus", Name: "ui.name.menus", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/posts", Name: "ui.name.posts", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/post/tags", Name: "ui.name.tags", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/post/categories", Name: "ui.name.categories", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/info-blocks", Name: "ui.name.info_blocks", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/templates", Name: "ui.name.templates", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/messages", Name: "ui.name.messages", Ico: template.HTML("<i class=\"material-icons\">mail_outline</i>")},
	}

	// Если передана функция перевода, переводим названия
	if tFunc != nil {
		for i := range routes {
			routes[i].Name = tFunc(routes[i].Name, nil)
		}
	}

	baseRoute := extractBaseRoute(currentRoute)
	for i := range routes {
		if routes[i].Path == baseRoute {
			routes[i].IsActive = true
		} else {
			routes[i].IsActive = false
		}
	}
	return routes
}

func extractBaseRoute(route string) string {
	parts := strings.Split(route, "/")
	baseRoute := ""
	for i, part := range parts {
		if i == 4 {
			break
		}
		if part == "" && i == 0 {
			continue
		}
		if !strings.HasPrefix(part, ":") {
			baseRoute += "/" + part
		} else {
			break
		}
	}
	return baseRoute
}
