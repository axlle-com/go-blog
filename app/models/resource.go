package models

import (
	"fmt"
	"os"
	"path/filepath"
)

type Resources struct {
	resources map[string]string
}

func NewResources() *Resources {
	return &Resources{
		resources: map[string]string{
			"posts":           "index",
			"post_categories": "post_categories",
			"info_blocks":     "info_blocks",
		},
	}
}

func (r *Resources) Resources() map[string]string {
	return r.resources
}

func (r *Resources) ResourceTemplate(name string) string {
	value, ok := r.resources[name]
	if !ok {
		return ""
	}

	fileName := filepath.Base(fmt.Sprintf("%s.gohtml", value))
	templatePath := filepath.Join("src/templates", fileName)

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return ""
	}

	return string(data)
}
