package menu

import (
	"html/template"
	"strings"
)

type Menu struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Ico      template.HTML `json:"ico"`
	IsActive bool          `json:"is_active"`
}

func NewMenu(currentRoute string) []Menu {
	var routes = []Menu{
		{Path: "/admin", Name: "Dashboard", Ico: template.HTML("<i data-feather=\"globe\"></i>")},
		{Path: "/admin/posts", Name: "Посты", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
	}

	for i := range routes {
		if extractBaseRoute(routes[i].Path) == currentRoute {
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
