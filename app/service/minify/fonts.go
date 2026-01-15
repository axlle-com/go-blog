package minify

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
)

// extractFontUrls извлекает все URL шрифтов из CSS контента
func extractFontUrls(cssContent string) []string {
	// Регулярное выражение для поиска url(...) в CSS
	// Поддерживает различные форматы: url(font.woff), url('font.woff'), url("font.woff"), url(fonts/font.woff)
	urlRegex := regexp.MustCompile(`url\(['"]?([^'")]+\.(woff2?|ttf|eot|otf))['"]?\)`)
	matches := urlRegex.FindAllStringSubmatch(cssContent, -1)

	var fontUrls []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		fontPath := strings.TrimSpace(match[1])

		// Пропускаем data: URLs и абсолютные пути
		if strings.HasPrefix(fontPath, "data:") || strings.HasPrefix(fontPath, "http://") || strings.HasPrefix(fontPath, "https://") {
			continue
		}

		// Извлекаем имя файла из пути
		fontName := filepath.Base(fontPath)
		if fontName != "" && !seen[fontName] {
			seen[fontName] = true
			fontUrls = append(fontUrls, fontName)
		}
	}

	return fontUrls
}

// findFontInResources ищет файл шрифта в src/resources/font (рекурсивно)
func findFontInResources(fontName string) (string, error) {
	fontPath := config.Config().DataFolder("resources", "font")

	// Используем filepath.Walk для поиска файла по имени (без учета регистра)
	var foundPath string
	err := filepath.Walk(fontPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.EqualFold(filepath.Base(path), fontName) {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if err == nil && foundPath != "" {
		return foundPath, nil
	}

	return "", os.ErrNotExist
}

// copyFontToStatic копирует шрифт в src/static/fonts
func copyFontToStatic(fontSourcePath string, fontName string) error {
	staticFontsDir := config.Config().StaticFolder("fonts")

	// Создаем директорию если её нет
	if err := os.MkdirAll(staticFontsDir, 0755); err != nil {
		return err
	}

	destPath := filepath.Join(staticFontsDir, fontName)

	// Читаем исходный файл
	sourceData, err := os.ReadFile(fontSourcePath)
	if err != nil {
		return err
	}

	// Записываем в назначение
	if err := os.WriteFile(destPath, sourceData, 0644); err != nil {
		return err
	}

	logger.Debugf("[web][copyFontToStatic] copied font: %s -> %s", fontSourcePath, destPath)
	return nil
}

// processFontsInCSS обрабатывает CSS контент: находит шрифты, копирует их и обновляет пути
func processFontsInCSS(cssContent string, cssFilePath string) (string, error) {
	fontUrls := extractFontUrls(cssContent)

	processedCSS := cssContent

	for _, fontName := range fontUrls {
		// Ищем шрифт в resources
		fontSourcePath, err := findFontInResources(fontName)
		if err != nil {
			logger.Debugf("[web][processFontsInCSS] font not found in resources: %s (referenced from %s)", fontName, cssFilePath)
			continue
		}

		// Копируем в static/fonts
		if err := copyFontToStatic(fontSourcePath, fontName); err != nil {
			logger.Errorf("[web][processFontsInCSS] failed to copy font %s: %v", fontName, err)
			continue
		}

		// Обновляем пути в CSS на /static/fonts/fontName
		// Заменяем различные варианты путей: fonts/font.woff, ../fonts/font.woff, font.woff и т.д.
		oldPattern := regexp.MustCompile(`url\(['"]?([^'")]*` + regexp.QuoteMeta(fontName) + `)['"]?\)`)
		newPath := "/static/fonts/" + fontName
		processedCSS = oldPattern.ReplaceAllStringFunc(processedCSS, func(match string) string {
			// Извлекаем текущий путь
			submatch := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`).FindStringSubmatch(match)
			if len(submatch) >= 2 {
				return `url(` + newPath + `)`
			}
			return match
		})
	}

	return processedCSS, nil
}
