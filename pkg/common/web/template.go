package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"path/filepath"
)

func InitTemplate(router *gin.Engine) {
	router.Static("/public", "./src/public")
	templates := loadTemplates("src/templates")
	router.SetHTMLTemplate(templates)
	//router.LoadHTMLGlob("templates/**/**/*")
}

func loadTemplates(templatesDir string) *template.Template {
	tmpl := template.New("")
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gohtml" {
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return tmpl
}
