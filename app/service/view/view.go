package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html/template"
	"io/fs"
	"net/url"
	"path"
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
	assetVersions    map[string]int64 // Кеш версий файлов: путь -> hash
	versionsMutex    sync.RWMutex
	dynamicTemplates map[string]string
	tmpl             *template.Template
	diskService      contract.DiskService
}

func NewView(config contract.Config, diskService contract.DiskService) *View {
	return &View{
		config:           config,
		assetVersions:    make(map[string]int64),
		dynamicTemplates: make(map[string]string),
		diskService:      diskService,
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

	templates := v.loadTemplatesFromFS(v.diskService.GetTemplatesFS(), "templates")
	v.tmpl = templates
	v.router.SetHTMLTemplate(templates)

	v.loadAssetVersionsFromFS(v.diskService.GetPublicFS(), "public")
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
	if v.router == nil {
		return fmt.Errorf("router is nil")
	}

	fullName := v.View(name)
	v.dynamicTemplates[fullName] = tmplStr

	baseTmpl := v.loadTemplatesFromFS(v.diskService.GetTemplatesFS(), "templates")

	for tName, tStr := range v.dynamicTemplates {
		newTmpl := baseTmpl.New(tName)
		if _, err := newTmpl.Parse(tStr); err != nil {
			logger.Errorf("[app][service][view][AddTemplateFromString] parse %q failed: %v", tName, err)
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

	return fmt.Sprintf("%s.%s", v.config.Layout(), name)
}

func (v *View) removeWhitespaceBetweenTags(s string) string {
	re := regexp.MustCompile(`>\s+<`)
	compactHTML := re.ReplaceAllString(s, "><")
	compactHTML = regexp.MustCompile(`[\n\r\t]+`).ReplaceAllString(compactHTML, " ")
	compactHTML = regexp.MustCompile(`\s+`).ReplaceAllString(compactHTML, " ")
	return strings.TrimSpace(compactHTML)
}

func (v *View) loadTemplatesFromFS(fsys fs.FS, templatesRoot string) *template.Template {
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
				logger.Errorf("[app][template][render] render template %q failed: %v", name, err)
				return template.HTML(fmt.Sprintf("<!-- render %q error: %v -->", name, err))
			}
			return template.HTML(buf.String())
		},
	}

	tmpl = tmpl.Funcs(funcMap)

	var files []string
	err := fs.WalkDir(fsys, templatesRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(p, ".gohtm") || strings.HasSuffix(p, ".gohtml") {
			files = append(files, p)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic("no templates found (embed) in " + templatesRoot)
	}

	if _, err := tmpl.ParseFS(fsys, files...); err != nil {
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

// loadAssetVersionsFromFS загружает версии (hash) всех статических JS/CSS файлов при старте
// В embed.FS нет ModTime, поэтому версия = crc32 от содержимого.
func (v *View) loadAssetVersionsFromFS(fsys fs.FS, publicRoot string) {
	v.versionsMutex.Lock()
	defer v.versionsMutex.Unlock()

	v.assetVersions = make(map[string]int64)

	// проверим, что директория существует
	if _, err := fs.Stat(fsys, publicRoot); err != nil {
		logger.Errorf("[template][loadAssetVersions] public directory not found in FS: %s (%v)", publicRoot, err)
		return
	}

	err := fs.WalkDir(fsys, publicRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(path.Ext(p))
		if ext != ".js" && ext != ".css" {
			return nil
		}

		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}

		sum := crc32.ChecksumIEEE(b) // uint32
		rel := strings.TrimPrefix(p, publicRoot+"/")
		assetPath := "/public/" + rel

		v.assetVersions[assetPath] = int64(sum)
		return nil
	})

	if err != nil {
		logger.Errorf("[template][loadAssetVersions] error walking public directory: %v", err)
	}

	logger.Infof("[template][loadAssetVersions] loaded %d asset versions", len(v.assetVersions))
}

// asset возвращает путь к файлу с версией на основе хеша содержимого
func (v *View) asset(p string) string {
	v.versionsMutex.RLock()
	defer v.versionsMutex.RUnlock()

	normalizedPath := normalizeAssetPath(p)

	version, exists := v.assetVersions[normalizedPath]
	if !exists {
		logger.Debugf("[template][asset] file not found in cache: %s", normalizedPath)
		return normalizedPath
	}

	return fmt.Sprintf("%s?v=%d", normalizedPath, version)
}

func normalizeAssetPath(p string) string {
	// URL пути должны быть с /
	p = strings.ReplaceAll(p, `\`, `/`)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	// path.Clean убирает // и ../
	p = path.Clean(p)
	return p
}

func (v *View) css(p string) template.HTML {
	href := v.asset(p)
	return template.HTML(
		`<link rel="preload" href="` + href + `" as="style" ` +
			`onload="this.onload=null;this.rel='stylesheet'">` +
			`<noscript><link rel="stylesheet" href="` + href + `"></noscript>`,
	)
}

func (v *View) js(p string) template.HTML {
	href := v.asset(p)
	return template.HTML(`<script src="` + href + `" defer></script>`)
}
