package minify

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type WebMinifier struct {
	config      contract.Config
	diskService contract.DiskService
	minifier    *minify.M
	mu          sync.RWMutex
	cache       map[string]string
}

func NewWebMinifier(cfg contract.Config, diskService contract.DiskService) contract.Minifier {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("application/javascript", js.Minify)

	return &WebMinifier{
		config:      cfg,
		diskService: diskService,
		minifier:    minifier,
		cache:       make(map[string]string),
	}
}

func (w *WebMinifier) Run() error {
	return w.cleanBundlesDir()
}

func (w *WebMinifier) Bundle(mediaType string, inputPaths []string) (string, error) {
	key := w.makeKey(mediaType, inputPaths)

	w.mu.RLock()
	if url, ok := w.cache[key]; ok {
		w.mu.RUnlock()
		return url, nil
	}
	w.mu.RUnlock()

	ext, err := extByMediaType(mediaType)
	if err != nil {
		return "", err
	}

	// public/bundles/<hash>.min.css
	rel := filepath.ToSlash(filepath.Join("public", "bundles", key+".min"+ext))
	outPath := w.config.DataFolder(filepath.FromSlash(rel))
	url := "/" + rel

	// 2) если файл уже есть на диске — не пересобираем
	if _, err := os.Stat(outPath); err == nil {
		w.mu.Lock()
		w.cache[key] = url
		w.mu.Unlock()
		return url, nil
	}

	// 3) собираем/минифицируем
	if err := w.mergeAndMinifyFiles(mediaType, inputPaths, outPath); err != nil {
		return "", err
	}

	w.mu.Lock()
	w.cache[key] = url
	w.mu.Unlock()

	return url, nil
}

func (w *WebMinifier) mergeAndMinifyFiles(mediaType string, inputPaths []string, outputPath string) error {
	var buffer bytes.Buffer
	var importsBuffer bytes.Buffer

	for _, inputPath := range inputPaths {
		input, err := w.diskService.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("[mergeAndMinifyFiles] read %s: %w", inputPath, err)
		}

		inputStr := string(input)

		// Если это CSS, обрабатываем шрифты
		if mediaType == "text/css" {
			// Извлекаем @import правила и удаляем их из содержимого
			imports, content := extractImports(inputStr)
			if len(imports) > 0 {
				importsBuffer.WriteString(imports)
				importsBuffer.WriteString("\n")
			}
			inputStr = content
		}

		buffer.WriteString(inputStr)
		buffer.WriteString("\n")
	}

	// Если есть @import правила, помещаем их в начало
	var finalBuffer bytes.Buffer
	if importsBuffer.Len() > 0 {
		finalBuffer.Write(importsBuffer.Bytes())
		finalBuffer.WriteString("\n")
	}
	finalBuffer.Write(buffer.Bytes())
	buffer = finalBuffer

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("[mergeAndMinifyFiles] mkdir %s: %w", filepath.Dir(outputPath), err)
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := output.Close(); cerr != nil {
			logger.Errorf("[mergeAndMinifyFiles] close %s: %v", outputPath, cerr)
		}
	}()

	if err := w.minifier.Minify(mediaType, output, &buffer); err != nil {
		return fmt.Errorf("[mergeAndMinifyFiles] minify %s: %w", outputPath, err)
	}

	return nil
}

func (w *WebMinifier) makeKey(mediaType string, inputPaths []string) string {
	// порядок важен => просто join в строку
	s := mediaType + "\n" + strings.Join(inputPaths, "\n")
	sum := crc32.ChecksumIEEE([]byte(s))
	// hex от crc32 — короткий и норм для имени файла
	return fmt.Sprintf("%08x", sum)
}

func (w *WebMinifier) cleanBundlesDir() error {
	// data/public/bundles
	dir := w.config.DataFolder(filepath.Join("public", "bundles"))

	// удалить всё содержимое (и саму папку, если есть)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove %s: %w", dir, err)
	}

	// создать заново
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	return nil
}

func extByMediaType(mediaType string) (string, error) {
	switch mediaType {
	case "text/css":
		return ".css", nil
	case "application/javascript":
		return ".js", nil
	default:
		return "", fmt.Errorf("unsupported mediaType: %s", mediaType)
	}
}

// extractImports извлекает все @import правила из CSS и возвращает их отдельно от остального содержимого
func extractImports(cssContent string) (imports string, content string) {
	// Регулярное выражение для поиска @import правил (поддерживает различные форматы)
	importRegex := regexp.MustCompile(`@import\s+(?:url\()?['"]?([^'"]+)['"]?\)?[^;]*;?`)

	var importsList []string
	var found bool

	// Находим все @import правила
	matches := importRegex.FindAllString(cssContent, -1)
	if len(matches) > 0 {
		found = true
		importsList = matches
	}

	// Удаляем @import правила из содержимого
	content = importRegex.ReplaceAllString(cssContent, "")
	content = strings.TrimSpace(content)

	// Формируем строку с @import правилами
	if found {
		imports = strings.Join(importsList, "\n")
	}

	return imports, content
}
