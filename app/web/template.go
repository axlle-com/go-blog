package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-gonic/gin"
)

var dynamicTemplates = make(map[string]string)

type Template struct {
	config        contract.Config
	router        *gin.Engine
	assetVersions map[string]int64 // Кеш версий файлов: путь -> timestamp
	versionsMutex sync.RWMutex
}

var (
	instance *Template
	once     sync.Once
)

func NewTemplate(router *gin.Engine) *Template {
	once.Do(func() {
		instance = &Template{
			config:        config.Config(),
			router:        router,
			assetVersions: make(map[string]int64),
		}
		instance.init()
		instance.loadAssetVersions()
	})
	return instance
}

func (t *Template) init() {
	templates := t.loadTemplates(t.config.SrcFolderBuilder("templates"))
	t.router.SetHTMLTemplate(templates)

	t.router.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/public/") || p == "/favicon.ico" {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	})

	t.router.StaticFile("/favicon.ico", "./"+t.config.SrcFolderBuilder("public/favicon.ico"))
	t.router.Static("/public", "./"+t.config.SrcFolderBuilder("public"))
	//router.LoadHTMLGlob("templates/**/**/*")
}

func (t *Template) ReLoad() {
	templates := t.loadTemplates(t.config.SrcFolderBuilder("templates"))
	t.router.SetHTMLTemplate(templates)
	t.loadAssetVersions()
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
	tmpl := template.New("")

	funcMap := template.FuncMap{
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
		"asset":      t.asset,
		"css":        t.css,
		"js":         t.js,
		"hasPrefix":  strings.HasPrefix,

		"render": func(name string, data any) template.HTML {
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
				logger.Errorf("[app][template][loadTemplates] render template %q failed: %v", name, err)
				return template.HTML(fmt.Sprintf("<!-- render %q error: %v -->", name, err))
			}
			return template.HTML(buf.String())
		},
	}

	tmpl = tmpl.Funcs(funcMap)

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gohtml" {
			if _, err = tmpl.ParseFiles(path); err != nil {
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

// loadAssetVersions загружает версии всех статических файлов при старте
func (t *Template) loadAssetVersions() {
	t.versionsMutex.Lock()
	defer t.versionsMutex.Unlock()

	// Очищаем старый кеш
	t.assetVersions = make(map[string]int64)

	publicDir := t.config.SrcFolderBuilder("public")
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		logger.Errorf("[template][loadAssetVersions] public directory not found: %s", publicDir)
		return
	}

	// Проходим по всем файлам в public директории
	err := filepath.Walk(publicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем только JS и CSS файлы
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".js" || ext == ".css" {
				// Получаем относительный путь от public директории
				relPath, err := filepath.Rel(publicDir, path)
				if err != nil {
					return err
				}

				// Нормализуем путь (заменяем обратные слеши на прямые)
				relPath = filepath.ToSlash(relPath)
				// Добавляем префикс /public/
				assetPath := "/public/" + relPath

				// Получаем время модификации файла
				modTime := info.ModTime().Unix()
				t.assetVersions[assetPath] = modTime
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorf("[template][loadAssetVersions] error walking public directory: %v", err)
	}

	logger.Infof("[template][loadAssetVersions] loaded %d asset versions", len(t.assetVersions))
}

// asset возвращает путь к файлу с версией на основе времени модификации
func (t *Template) asset(path string) string {
	t.versionsMutex.RLock()
	defer t.versionsMutex.RUnlock()

	// Нормализуем путь
	normalizedPath := filepath.ToSlash(path)
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// Получаем версию из кеша
	version, exists := t.assetVersions[normalizedPath]
	if !exists {
		// Если файл не найден в кеше, возвращаем путь без версии
		logger.Debugf("[template][asset] file not found in cache: %s", normalizedPath)
		return normalizedPath
	}

	// Добавляем версию как query параметр
	return fmt.Sprintf("%s?v=%d", normalizedPath, version)
}

func (t *Template) css(path string) template.HTML {
	href := t.asset(path)

	return template.HTML(
		`<link rel="preload" href="` + href + `" as="style" ` +
			`onload="this.onload=null;this.rel='stylesheet'">` +
			`<noscript><link rel="stylesheet" href="` + href + `"></noscript>`,
	)
}

func (t *Template) js(path string) template.HTML {
	href := t.asset(path)

	return template.HTML(
		`<script src="` + href + `" defer></script>`,
	)
}
