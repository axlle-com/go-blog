package models

import (
	"fmt"
	"os"
	"path/filepath"
)

type ResourceMap struct {
	resources map[string]string
}

func NewResource() *ResourceMap {
	return &ResourceMap{
		resources: map[string]string{
			"posts":           "index",
			"post_categories": "post_categories",
			"info_blocks":     "info_blocks",
		},
	}
}

func (r *ResourceMap) Resources() map[string]string {
	return r.resources
}

func (r *ResourceMap) ResourceTemplate(name string) string {
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
