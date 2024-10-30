package web

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func InitTemplate(router *gin.Engine) {
	cfg := config.Config()
	router.StaticFile("/favicon.ico", "./"+cfg.SrcFolderBuilder("public/favicon.ico"))
	router.Static("/public", "./"+cfg.SrcFolderBuilder("public"))
	templates := LoadTemplates(cfg.SrcFolderBuilder("templates"))
	router.SetHTMLTemplate(templates)
	//router.LoadHTMLGlob("templates/**/**/*")
}

func LoadTemplates(templatesDir string) *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"add":     add,
		"sub":     sub,
		"mul":     mul,
		"date":    date,
		"query":   query,
		"ptrStr":  ptrStr,
		"ptrUint": ptrUint,
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
func query(url url.Values, string string) string {
	param, ok := url[string]
	if ok {
		// TODO
		return param[0]
	}
	return ""
}
func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func ptrUint(i *uint) uint {
	if i == nil {
		return 0
	}
	return *i
}
