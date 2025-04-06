package web

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/axlle-com/blog/pkg/app/config"
	"github.com/gin-gonic/gin"
)

type Template struct {
	config contracts.Config
	router *gin.Engine
}

var (
	instance *Template
	once     sync.Once
)

func NewTemplate(router *gin.Engine) *Template {
	once.Do(func() {
		instance = &Template{
			config: config.Config(),
			router: router,
		}
		instance.init()
	})
	return instance
}

func (t *Template) init() {
	t.router.StaticFile("/favicon.ico", "./"+t.config.SrcFolderBuilder("public/favicon.ico"))
	t.router.Static("/public", "./"+t.config.SrcFolderBuilder("public"))
	templates := t.loadTemplates(t.config.SrcFolderBuilder("templates"))
	t.router.SetHTMLTemplate(templates)
	//router.LoadHTMLGlob("templates/**/**/*")
}

func (t *Template) ReLoad() {
	templates := t.loadTemplates(t.config.SrcFolderBuilder("templates"))
	t.router.SetHTMLTemplate(templates)
}

func (t *Template) loadTemplates(templatesDir string) *template.Template {
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

func query(url url.Values, key string) string {
	param, ok := url[key]
	if ok && len(param) > 0 {
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
