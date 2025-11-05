package web

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
)

var (
	safelist         = []string{"show", "active", "is-open", "is-active"}
	safelistPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^modal`),
		regexp.MustCompile(`^collapse`),
		regexp.MustCompile(`^dropdown`),
		regexp.MustCompile(`^owl-`),
		regexp.MustCompile(`^swiper-`),
		regexp.MustCompile(`^slick-`),
		regexp.MustCompile(`^js-`),
	}
	whitelistedTags = []string{
		"html", "body", "a", "p", "ul", "ol", "li", "img", "button",
		"h1", "h2", "h3", "h4", "h5", "h6", "input", "textarea", "label",
		"section", "header", "footer", "nav", "main", "small", "strong", "em",
	}
	reClass      = regexp.MustCompile(`class\s*=\s*"(.*?)"|class\s*=\s*'(.*?)'`)
	reWord       = regexp.MustCompile(`[A-Za-z0-9_\-:/.%]+`)
	reJSClassAdd = regexp.MustCompile(`(?i)(add|remove|toggle)(Class|)\(([^)]*)\)`)
	reID         = regexp.MustCompile(`id\s*=\s*"(.*?)"|id\s*=\s*'(.*?)'`)

	reDotClass = regexp.MustCompile(`\.[A-Za-z0-9_\-:/.%]+`)
	reHashID   = regexp.MustCompile(`\#[A-Za-z0-9_\-:/.%]+`)
)

// упрощённый recursive-walk вместо ** в Glob
func collectFiles(root string, ext ...string) []string {
	extset := map[string]struct{}{}
	for _, e := range ext {
		extset[strings.ToLower(e)] = struct{}{}
	}

	var out []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if len(extset) == 0 {
			out = append(out, path)
			return nil
		}
		if _, ok := extset[strings.ToLower(filepath.Ext(path))]; ok {
			out = append(out, path)
		}
		return nil
	})
	return out
}

func startCompress(cfg contract.Config, inputCSS string) string {
	used := make(map[string]struct{}) // ".class" / "#id" / "tag"

	// где искать разметку
	tplRoot := filepath.Join("templates", "front", cfg.Layout())
	jsRoot := cfg.PublicFolderBuilder(cfg.Layout(), "js")

	var contentFiles []string
	contentFiles = append(contentFiles, collectFiles(tplRoot, ".gohtml", ".html")...)
	contentFiles = append(contentFiles, collectFiles(jsRoot, ".js")...)

	for _, f := range contentFiles {
		b, err := os.ReadFile(f)
		if err != nil {
			logger.Debugf("skip %s: %v", f, err)
			continue
		}
		t := string(b)

		for _, m := range reClass.FindAllStringSubmatch(t, -1) {
			val := m[1]
			if val == "" {
				val = m[2]
			}
			for _, w := range reWord.FindAllString(val, -1) {
				used["."+w] = struct{}{}
			}
		}
		for _, m := range reJSClassAdd.FindAllStringSubmatch(t, -1) {
			args := m[3]
			for _, w := range reWord.FindAllString(args, -1) {
				used["."+w] = struct{}{}
			}
		}
		for _, m := range reID.FindAllStringSubmatch(t, -1) {
			val := m[1]
			if val == "" {
				val = m[2]
			}
			for _, w := range reWord.FindAllString(val, -1) {
				used["#"+w] = struct{}{}
			}
		}
	}

	for _, tag := range whitelistedTags {
		used[tag] = struct{}{}
	}
	for _, s := range safelist {
		used["."+s] = struct{}{}
	}

	in, err := os.ReadFile(inputCSS)
	if err != nil {
		logger.Errorf("read css: %v", err)
		return ""
	}

	blocks := splitRules(string(in))

	var out bytes.Buffer
	bw := bufio.NewWriter(&out)
	kept := 0

	for _, b := range blocks {
		selector, body := splitSelectorBody(b)
		if selector == "" || body == "" {
			continue
		}

		// @-правила
		if strings.HasPrefix(selector, "@") {
			if keepAtRule(selector, body, used) {
				fmt.Fprintf(bw, "%s{%s}\n", selector, body)
				kept++
			}
			continue
		}

		if keepRule(selector, used) {
			fmt.Fprintf(bw, "%s{%s}\n", selector, body)
			kept++
		}
	}

	_ = bw.Flush()

	outputCSS := cfg.PublicFolderBuilder(cfg.Layout(), "css", "app.min.css")
	if err := os.WriteFile(outputCSS, out.Bytes(), 0644); err != nil {
		logger.Errorf("write out: %v", err)
		return ""
	}

	logger.Infof("Purged: kept %d rules → %s", kept, outputCSS)

	return outputCSS
}

func splitRules(css string) []string {
	var res []string
	var cur strings.Builder
	depth := 0
	for _, r := range css {
		cur.WriteRune(r)
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				res = append(res, cur.String())
				cur.Reset()
			}
		}
	}
	return res
}

func splitSelectorBody(rule string) (selector, body string) {
	i := strings.IndexByte(rule, '{')
	j := strings.LastIndexByte(rule, '}')
	if i < 0 || j < 0 || j <= i {
		return "", ""
	}
	selector = strings.TrimSpace(rule[:i])
	body = strings.TrimSpace(rule[i+1 : j])
	return
}

func keepRule(selector string, used map[string]struct{}) bool {
	for _, part := range strings.Split(selector, ",") {
		p := strings.TrimSpace(part)

		for _, tag := range whitelistedTags {
			if p == tag || strings.HasPrefix(p, tag+" ") || strings.HasPrefix(p, tag+".") || strings.HasPrefix(p, tag+"#") {
				return true
			}
		}
		for _, m := range reDotClass.FindAllString(p, -1) {
			if _, ok := used[m]; ok {
				return true
			}
			for _, rx := range safelistPatterns {
				if rx.MatchString(strings.TrimPrefix(m, ".")) {
					return true
				}
			}
		}
		for _, m := range reHashID.FindAllString(p, -1) {
			if _, ok := used[m]; ok {
				return true
			}
		}
	}
	return false
}

// простая логика для @-правил
func keepAtRule(selector, body string, used map[string]struct{}) bool {
	sel := strings.ToLower(selector)
	switch {
	case strings.HasPrefix(sel, "@font-face"),
		strings.HasPrefix(sel, "@keyframes"),
		strings.HasPrefix(sel, "@supports"):
		return true
	case strings.HasPrefix(sel, "@media"):
		// разберём вложенные правила и оставим блок, если внутри есть хоть одно нужное
		child := splitRules(selector + "{" + body + "}")
		if len(child) == 0 {
			return false
		}
		// первый элемент — сам @media, дальше — дети
		for _, r := range child[1:] {
			s, b := splitSelectorBody(r)
			if s == "" || b == "" {
				continue
			}
			if strings.HasPrefix(s, "@") {
				if keepAtRule(s, b, used) {
					return true
				}
				continue
			}
			if keepRule(s, used) {
				return true
			}
		}
		return false
	default:
		// на всякий случай — сохраняем неизвестные @-блоки
		return true
	}
}
