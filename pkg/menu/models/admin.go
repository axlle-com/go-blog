package models

import (
	"html/template"
	"strings"
)

type AdminMenu struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Ico      template.HTML `json:"ico"`
	IsActive bool          `json:"is_active"`
}

func NewMenu(currentRoute string) []AdminMenu {
	var routes = []AdminMenu{
		{Path: "/admin/", Name: "Dashboard", Ico: template.HTML("<i data-feather=\"globe\"></i>")},
		{Path: "/admin/posts", Name: "Посты", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
		{Path: "/admin/categories", Name: "Категории", Ico: template.HTML("<i class=\"material-icons\">list_alt</i>")},
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
		if i == 3 {
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
