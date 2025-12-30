package models

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/axlle-com/blog/src" // @todo избавиться
)

type Resources struct {
	resources map[string]string
	themes    map[string]string
}

func NewResources() *Resources {
	return &Resources{
		resources: map[string]string{
			"posts":           "index",
			"post_categories": "post_categories",
			"post_tags":       "post_tags",
			"info_blocks":     "info_blocks",
			"menus":           "menus",
		},
		themes: map[string]string{
			"default": "default",
			"spring":  "spring",
		},
	}
}

func (r *Resources) Resources() map[string]string { return r.resources }
func (r *Resources) Themes() map[string]string    { return r.themes }

// ResourceTemplate возвращает исходник шаблона из embed (templates/**/<file>.gohtml).
// Ищет файл по базовому имени (например index.gohtml) по всему дереву templates.
// Если найдены несколько — берём первый.
func (r *Resources) ResourceTemplate(name string) string {
	value, ok := r.resources[name]
	if !ok {
		return ""
	}

	wantGohtml := fmt.Sprintf("%s.gohtml", value)

	var foundPath string

	_ = fs.WalkDir(src.TemplatesFS, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		base := path.Base(p)
		if base == wantGohtml {
			foundPath = p
			return fs.SkipAll
		}
		return nil
	})

	if foundPath == "" {
		return ""
	}

	b, err := fs.ReadFile(src.TemplatesFS, foundPath)
	if err != nil {
		return ""
	}

	// нормализуем переносы/нулевые байты не надо; вернём как есть
	return strings.TrimPrefix(string(b), "\uFEFF") // на случай BOM
}
