package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path"
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
	assetVersions    map[string]int64 // Кеш версий файлов: путь -> hash
	versionsMutex    sync.RWMutex
	dynamicTemplates map[string]string
	tmpl             *template.Template
	diskService      contract.DiskService
	minifier         contract.Minifier
}

func NewView(config contract.Config, diskService contract.DiskService, minifier contract.Minifier) *View {
	return &View{
		config:           config,
		assetVersions:    make(map[string]int64),
		dynamicTemplates: make(map[string]string),
		diskService:      diskService,
		minifier:         minifier,
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

	templates := v.loadTemplatesCombined(
		v.diskService.GetTemplatesFS(),
		"templates",
		v.config.DataFolder("templates"),
	)

	v.tmpl = templates
	v.router.SetHTMLTemplate(templates)

	v.loadAssetVersionsFromFS(v.diskService.GetStaticFS(), "static")
}

func (v *View) RenderToString(name string, data any) (string, error) {
	if v.tmpl == nil {
		return "", fmt.Errorf("template engine is not initialized")
	}

	name = v.ViewStatic(name)

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

	fullName := v.ViewStatic(name)
	v.dynamicTemplates[fullName] = tmplStr

	baseTmpl := v.loadTemplatesCombined(
		v.diskService.GetTemplatesFS(),
		"templates",
		v.config.DataFolder("templates"),
	)

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

func (v *View) ViewStatic(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) >= 3 {
		return name
	}

	if name == "" {
		name = "default"
	}

	return fmt.Sprintf("%s.%s", v.config.Layout(), name)
}

func (v *View) ViewResource(resource contract.Resource) string {
	if resource == nil {
		return ""
	}

	if resource.GetTemplateName() != "" {
		return resource.GetTemplateName()
	}

	return fmt.Sprintf("%s.%s.default", v.config.Layout(), resource.GetName())
}

func (v *View) removeWhitespaceBetweenTags(s string) string {
	re := regexp.MustCompile(`>\s+<`)
	compactHTML := re.ReplaceAllString(s, "><")
	compactHTML = regexp.MustCompile(`[\n\r\t]+`).ReplaceAllString(compactHTML, " ")
	compactHTML = regexp.MustCompile(`\s+`).ReplaceAllString(compactHTML, " ")

	return strings.TrimSpace(compactHTML)
}

func (v *View) loadTemplatesCombined(embedFS fs.FS, embedRoot string, diskRoot string) *template.Template {
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
		"hasPrefix":  strings.HasPrefix,

		"concatSlice":    concatSlice,
		"stringsToSlice": stringsToSlice,

		"js":             v.js,
		"css":            v.css,
		"asset":          v.asset,
		"bundleJs":       v.bundleJs,
		"bundleCss":      v.bundleCss,
		"bundleCssSlice": v.bundleCssSlice,
		"render": func(name string, data any) template.HTML {
			var buf bytes.Buffer
			name = v.ViewStatic(name)
			if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
				logger.Errorf("[app][template][render] render template %q failed: %v", name, err)
				return template.HTML(fmt.Sprintf("<!-- render %q error: %v -->", name, err))
			}

			return template.HTML(buf.String())
		},
	}

	tmpl = tmpl.Funcs(funcMap)

	embedFiles, err := listTemplateFilesFS(embedFS, embedRoot)
	if err != nil {
		panic(err)
	}

	if len(embedFiles) != 0 {
		if _, err := tmpl.ParseFS(embedFS, embedFiles...); err != nil {
			panic(err)
		}
	}

	if diskRoot != "" {
		if st, err := os.Stat(diskRoot); err == nil && st.IsDir() {
			diskFiles, err := listTemplateFilesDisk(diskRoot)
			if err != nil {
				panic(err)
			}

			if len(diskFiles) > 0 {
				if _, err := tmpl.ParseFiles(diskFiles...); err != nil {
					panic(err)
				}

				logger.Infof("[view][loadTemplatesCombined] loaded %d disk templates from %s", len(diskFiles), diskRoot)
			}
		} else {
			logger.Infof("[view][loadTemplatesCombined] disk templates dir not found: %s", diskRoot)
		}
	}

	return tmpl
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
	return template.HTML(`<script src="` + v.asset(p) + `" defer></script>`)
}

func (v *View) bundleJs(paths ...string) template.HTML {
	if v.config.IsLocal() {
		var b strings.Builder
		for _, p := range paths {
			src := v.asset(p)
			b.WriteString(`<script src="`)
			b.WriteString(src)
			b.WriteString(`" defer></script>` + "\n")
		}

		return template.HTML(b.String())
	}

	href, err := v.minifier.Bundle("application/javascript", paths)
	if err != nil {
		logger.Errorf("[template][bundleJs] %v", err)

		return template.HTML(fmt.Sprintf("<!-- bundleJs error: %v -->", err))
	}

	return template.HTML(`<script src="` + href + `" defer></script>`)
}

func (v *View) bundleCss(paths ...string) template.HTML {
	if v.config.IsLocal() {
		var b strings.Builder
		for _, p := range paths {
			href := v.asset(p)
			b.WriteString(`<link rel="stylesheet" href="`)
			b.WriteString(href)
			b.WriteString(`">` + "\n")
		}

		return template.HTML(b.String())
	}

	href, err := v.minifier.Bundle("text/css", paths)
	if err != nil {
		logger.Errorf("[template][bundleCss] %v", err)

		return template.HTML(fmt.Sprintf("<!-- bundleCss error: %v -->", err))
	}

	return template.HTML(
		`<link rel="preload" href="` + href + `" as="style" ` +
			`onload="this.onload=null;this.rel='stylesheet'">` +
			`<noscript><link rel="stylesheet" href="` + href + `"></noscript>`,
	)
}

func (v *View) loadAssetVersionsFromFS(fsys fs.FS, staticRoot string) {
	v.versionsMutex.Lock()
	defer v.versionsMutex.Unlock()

	v.assetVersions = make(map[string]int64)

	// проверим, что директория существует
	if _, err := fs.Stat(fsys, staticRoot); err != nil {
		logger.Errorf("[template][loadAssetVersions] static directory not found in FS: %s (%v)", staticRoot, err)
		return
	}

	err := fs.WalkDir(fsys, staticRoot, func(p string, d fs.DirEntry, err error) error {
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
		rel := strings.TrimPrefix(p, staticRoot+"/")
		assetPath := "/static/" + rel

		v.assetVersions[assetPath] = int64(sum)
		return nil
	})

	if err != nil {
		logger.Errorf("[template][loadAssetVersions] error walking static directory: %v", err)
	}

	logger.Infof("[template][loadAssetVersions] loaded %d asset versions", len(v.assetVersions))
}

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

func (v *View) render(tmpl *template.Template, name string, data any) template.HTML {
	var buf bytes.Buffer
	name = v.ViewStatic(name)

	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		logger.Errorf("[app][template][render] render template %q failed: %v", name, err)

		return template.HTML(fmt.Sprintf("<!-- render %q error: %v -->", name, err))
	}

	return template.HTML(buf.String())
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

func listTemplateFilesFS(fsys fs.FS, root string) ([]string, error) {
	var files []string
	err := fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
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
	return files, err
}

func listTemplateFilesDisk(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		low := strings.ToLower(p)
		if strings.HasSuffix(low, ".gohtm") || strings.HasSuffix(low, ".gohtml") {
			files = append(files, p)
		}

		return nil
	})

	return files, err
}

func stringsToSlice(items ...string) []string { return items }

func concatSlice(a, b []string) []string {
	out := make([]string, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}

func (v *View) bundleCssSlice(paths []string) template.HTML {
	return v.bundleCss(paths...)
}
