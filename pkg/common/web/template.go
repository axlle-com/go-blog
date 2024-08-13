package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func InitTemplate(router *gin.Engine) {
	router.Static("/favicon.ico", "./src/public/favicon.ico")
	router.Static("/public", "./src/public")
	templates := loadTemplates("src/templates")
	router.SetHTMLTemplate(templates)
	//router.LoadHTMLGlob("templates/**/**/*")
}

func loadTemplates(templatesDir string) *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"add":   add,
		"sub":   sub,
		"mul":   mul,
		"date":  date,
		"query": query,
	})
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

func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }
func mul(x, y int) int { return x * y }
func date(timePtr *time.Time) string {
	if timePtr == nil {
		return ""
	}
	return timePtr.Format("02.01.2006 15:04:05")
}
func query(u url.Values, s string) string {
	param, ok := u[s]
	if ok {
		// TODO
		return param[0]
	}
	return ""
}
