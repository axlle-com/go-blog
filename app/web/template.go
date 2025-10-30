package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-gonic/gin"
)

var dynamicTemplates = make(map[string]string)

type Template struct {
	config contract.Config
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

func (t *Template) AddTemplateFromString(name, tmplStr string) error {
	dynamicTemplates[name] = tmplStr

	baseTmpl := t.loadTemplates(t.config.SrcFolderBuilder("templates"))

	for name, tmplStr := range dynamicTemplates {
		newTmpl := baseTmpl.New(name)
		if _, err := newTmpl.Parse(tmplStr); err != nil {
			logger.Error(err)
			continue
		}
	}

	t.router.SetHTMLTemplate(baseTmpl)
	return nil
}

func (t *Template) loadTemplates(templatesDir string) *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"add":        add,
		"sub":        sub,
		"mul":        mul,
		"dict":       dict,
		"date":       date,
		"query":      query,
		"ptrStr":     ptrStr,
		"ptrUint":    ptrUint,
		"json":       jsonFunc,
		"collapseID": collapseID,
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

func jsonFunc(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func collapseID(prefix string, id uint) string {
	return fmt.Sprintf("%s-%d", prefix, id)
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict: odd args")
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %d not string", i)
		}
		m[k] = values[i+1]
	}
	return m, nil
}
