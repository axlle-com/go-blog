package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-gonic/gin"
)

type View struct {
	config           contract.Config
	router           *gin.Engine
	assetVersions    map[string]int64 // Кеш версий файлов: путь -> timestamp
	versionsMutex    sync.RWMutex
	dynamicTemplates map[string]string
	tmpl             *template.Template
}

func NewView(config contract.Config) *View {
	return &View{
		config:        config,
		assetVersions: make(map[string]int64),
	}
}

func (v *View) SetRouter(router *gin.Engine) {
	if router == nil {
		logger.Fatal("[app][service][view][SetRouter] router is nil")
		return
	}

	v.router = router
}

func (v *View) Load() {
	if v.router == nil {
		logger.Error("[app][service][view][Load] router is nil")
		return
	}

	templates := v.loadTemplates(v.config.SrcFolderBuilder("templates"))
	v.tmpl = templates
	v.router.SetHTMLTemplate(templates)
	v.loadAssetVersions()
}

func (v *View) SetStatic() {
	if v.router == nil {
		logger.Error("[app][service][view][SetStatic] router is nil")
		return
	}

	v.router.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/public/") || p == "/favicon.ico" {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	})

	v.router.StaticFile("/favicon.ico", "./"+v.config.SrcFolderBuilder("public/favicon.ico"))
	v.router.Static("/public", "./"+v.config.SrcFolderBuilder("public"))
	//router.LoadHTMLGlob("templates/**/**/*")
}

func (v *View) RenderToString(name string, data any) (string, error) {
	if v.tmpl == nil {
		return "", fmt.Errorf("template engine is not initialized")
	}

	name = v.View(name)

	var buf bytes.Buffer
	if err := v.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}

	return v.removeWhitespaceBetweenTags(buf.String()), nil
}

func (v *View) AddTemplateFromString(name, tmplStr string) error {
	v.dynamicTemplates[name] = tmplStr

	baseTmpl := v.loadTemplates(v.config.SrcFolderBuilder("templates"))

	for name, tmplStr := range v.dynamicTemplates {
		newTmpl := baseTmpl.New(name)
		if _, err := newTmpl.Parse(tmplStr); err != nil {
			logger.Error(err)
			continue
		}
	}

	v.tmpl = baseTmpl
	v.router.SetHTMLTemplate(baseTmpl)
	return nil
}

func (v *View) View(name string) string {
	if name == "" {
		name = "default"
	}

	layout := v.config.Layout()
	if layout == "" {
		layout = "default"
	}

	return fmt.Sprintf("%s.%s", layout, name)
}

func (v *View) ViewStatic(name string) string {
	if name == "" {
		name = "index"
	}

	return fmt.Sprintf("%s.%s", v.config.Layout(), name)
}

func (v *View) removeWhitespaceBetweenTags(s string) string {
	// Удаляем пробелы между тегами
	re := regexp.MustCompile(`>\s+<`)
	compactHTML := re.ReplaceAllString(s, "><")
	// Удаляем все переносы строк и табы, заменяя их на пробелы
	compactHTML = regexp.MustCompile(`[\n\r\t]+`).ReplaceAllString(compactHTML, " ")
	// Удаляем множественные пробелы
	compactHTML = regexp.MustCompile(`\s+`).ReplaceAllString(compactHTML, " ")
	return strings.TrimSpace(compactHTML)
}

func (v *View) loadTemplates(templatesDir string) *template.Template {
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
		"asset":      v.asset,
		"css":        v.css,
		"js":         v.js,
		"hasPrefix":  strings.HasPrefix,

		"render": func(name string, data any) template.HTML {
			var buf bytes.Buffer
			name = v.View(name)
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

func splitFirst(s string) (first, next string) {
	parts := strings.SplitN(s, ".", 2)
	first = parts[0]
	if len(parts) > 1 {
		next = parts[1]
	}
	return
}

// loadAssetVersions загружает версии всех статических файлов при старте
func (v *View) loadAssetVersions() {
	v.versionsMutex.Lock()
	defer v.versionsMutex.Unlock()

	// Очищаем старый кеш
	v.assetVersions = make(map[string]int64)

	publicDir := v.config.SrcFolderBuilder("public")
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
				v.assetVersions[assetPath] = modTime
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorf("[template][loadAssetVersions] error walking public directory: %v", err)
	}

	logger.Infof("[template][loadAssetVersions] loaded %d asset versions", len(v.assetVersions))
}

// asset возвращает путь к файлу с версией на основе времени модификации
func (v *View) asset(path string) string {
	v.versionsMutex.RLock()
	defer v.versionsMutex.RUnlock()

	// Нормализуем путь
	normalizedPath := filepath.ToSlash(path)
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// Получаем версию из кеша
	version, exists := v.assetVersions[normalizedPath]
	if !exists {
		// Если файл не найден в кеше, возвращаем путь без версии
		logger.Debugf("[template][asset] file not found in cache: %s", normalizedPath)
		return normalizedPath
	}

	// Добавляем версию как query параметр
	return fmt.Sprintf("%s?v=%d", normalizedPath, version)
}

func (v *View) css(path string) template.HTML {
	href := v.asset(path)

	return template.HTML(
		`<link rel="preload" href="` + href + `" as="style" ` +
			`onload="this.onload=null;this.rel='stylesheet'">` +
			`<noscript><link rel="stylesheet" href="` + href + `"></noscript>`,
	)
}

func (v *View) js(path string) template.HTML {
	href := v.asset(path)

	return template.HTML(
		`<script src="` + href + `" defer></script>`,
	)
}
