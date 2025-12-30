package minify

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type WebMinifier struct {
	cfg contract.Config
	m   *minify.M
}

func NewWebMinifier(cfg contract.Config) *WebMinifier {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	return &WebMinifier{
		cfg: cfg,
		m:   m,
	}
}

func (w *WebMinifier) Run() error {
	resourcesDir := w.cfg.DataFolder("resources")
	if st, err := os.Stat(resourcesDir); err != nil {
		if os.IsNotExist(err) {
			logger.Infof("[web][Minify] resources dir not found (%s); skipping", resourcesDir)
			return nil
		}
		return fmt.Errorf("stat resources dir %q: %w", resourcesDir, err)
	} else if !st.IsDir() {
		logger.Infof("[web][Minify] resources path is not a directory (%s); skipping", resourcesDir)
		return nil
	}

	if err := w.minifyAdmin(); err != nil {
		return fmt.Errorf("minify admin: %w", err)
	}

	if w.cfg.Layout() == "spring" {
		return w.minifySpring()
	}

	return w.minifyFront()
}

func (w *WebMinifier) minifyAdmin() error {
	adminCSS := []string{
		"resources/font/inter/inter.min.css",
		"resources/font/play/play.css",
		"resources/plugins/material-design-icons-iconfont/material-design-icons.min.css",
		"resources/plugins/fontawesome-free/css/all.min.css",
		"resources/plugins/simplebar/simplebar.min.css",
		"resources/plugins/summernote/summernote-bs4.css",
		"resources/plugins/select2/css/select2.min.css",
		"resources/plugins/flatpickr/flatpickr.min.css",
		"resources/plugins/noty/noty.css",
		"resources/plugins/noty/themes/relax.css",
		"resources/plugins/fancybox/fancybox.css",
		"resources/admin/css/style.css",
		"resources/admin/css/sidebar-dark.min.css",
		"resources/plugins/sweetalert2/sweetalert2.min.css",
	}
	adminJS := []string{
		"resources/plugins/jquery-3-6-0/jquery.min.js",
		"resources/plugins/bootstrap-4-6-1/js/bootstrap.bundle.js",
		"resources/plugins/simplebar/simplebar.min.js",
		"resources/plugins/feather-icons/feather.min.js",
		"resources/plugins/summernote/summernote-bs4.min.js",
		"resources/plugins/select2/js/select2.full.js",
		"resources/plugins/select2/js/i18n/ru.js",
		"resources/plugins/flatpickr/flatpickr.js",
		"resources/plugins/flatpickr/l10n/ru.js",
		"resources/plugins/noty/noty.js",
		"resources/plugins/inputmask/jquery.inputmask.js",
		"resources/plugins/fancybox/fancybox.umd.js",
		"resources/plugins/sweetalert2/sweetalert2.all.min.js",
		"resources/plugins/sortablejs/Sortable.min.js",
		"resources/plugins/js/script.min.js",
		"resources/plugins/date-format/date-format.js",
	}

	if err := w.mergeAndMinifyFiles("text/css", adminCSS, w.cfg.PublicFolderBuilder("admin/app.css")); err != nil {
		return err
	}

	if err := w.mergeAndMinifyFiles("application/javascript", adminJS, w.cfg.PublicFolderBuilder("admin/app.js")); err != nil {
		return err
	}

	return nil
}

func (w *WebMinifier) minifyFront() error {
	CSS := []string{
		"resources/plugins/bootstrap-5.0.2-dist/css/bootstrap.min.css",
	}
	JS := []string{
		"resources/plugins/jquery-3-6-0/jquery.min.js",
		"resources/plugins/bootstrap-5.0.2-dist/js/bootstrap.min.js",
		"public/admin/glob.js",
	}

	if err := w.mergeAndMinifyFiles("text/css", CSS, w.cfg.PublicFolderBuilder("app.css")); err != nil {
		return err
	}

	if err := w.mergeAndMinifyFiles("application/javascript", JS, w.cfg.PublicFolderBuilder("app.js")); err != nil {
		return err
	}

	return nil
}

func (w *WebMinifier) minifySpring() error {
	CSS := []string{
		"public/spring/css/font.css",
		"resources/spring/css/bootstrap.css",
		"resources/spring/css/font-awesome.css",
		"resources/spring/css/themify-icons.css",
		"resources/spring/css/linear-icons.css",
		"resources/spring/css/animate.css",
		"resources/spring/css/owl.css",
		"resources/spring/css/jquery.fancybox.css",
		"resources/spring/css/style.css",
		"resources/spring/css/responsive.css",
		"resources/plugins/noty/noty.css",
		"resources/plugins/noty/themes/relax.css",
		"public/spring/css/common.css",
	}
	JS := []string{
		"resources/spring/js/jquery.js",
		"resources/spring/js/bootstrap.min.js",
		"resources/spring/js/appear.js",
		"resources/spring/js/pagenav.js",
		"resources/spring/js/jquery.scrollTo.js",
		"resources/spring/js/jquery.fancybox.pack.js",
		"resources/spring/js/owl.js",
		"resources/spring/js/wow.js",
		"resources/spring/js/validate.js",
		"resources/spring/js/script.js",
		"resources/plugins/noty/noty.js",
		"public/admin/glob.js",
	}

	if err := w.mergeAndMinifyFiles("text/css", CSS, w.cfg.PublicFolderBuilder("spring/css/app.css")); err != nil {
		return err
	}

	if err := w.mergeAndMinifyFiles("application/javascript", JS, w.cfg.PublicFolderBuilder("spring/js/app.js")); err != nil {
		return err
	}

	return nil
}

func (w *WebMinifier) mergeAndMinifyFiles(mediaType string, inputPaths []string, outputPath string) error {
	var buffer bytes.Buffer
	var importsBuffer bytes.Buffer

	for _, inputPath := range inputPaths {
		if strings.HasPrefix(inputPath, "resources") {
			inputPath = w.cfg.DataFolder(inputPath)
		}

		if strings.HasPrefix(inputPath, "public") {
			inputPath = w.cfg.SrcFolderBuilder(inputPath)
		}

		input, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("[mergeAndMinifyFiles] read %s: %w", inputPath, err)
		}

		inputStr := string(input)

		// Если это CSS, обрабатываем шрифты
		if mediaType == "text/css" {
			// Находим шрифты, копируем их в public/fonts и обновляем пути
			processedCSS, err := processFontsInCSS(inputStr, inputPath)
			if err != nil {
				logger.Errorf("[mergeAndMinifyFiles] error processing fonts in %s: %v", inputPath, err)
				// Продолжаем с оригинальным CSS в случае ошибки
				inputStr = string(input)
			} else {
				inputStr = processedCSS
			}

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

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := output.Close(); cerr != nil {
			logger.Errorf("[mergeAndMinifyFiles] close %s: %v", outputPath, cerr)
		}
	}()

	if err := w.m.Minify(mediaType, output, &buffer); err != nil {
		return fmt.Errorf("[mergeAndMinifyFiles] minify %s: %w", outputPath, err)
	}

	return nil
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
