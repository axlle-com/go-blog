package web

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"io/ioutil"
	"log"
	"os"
)

func InitMinify() {
	adminCSS := []string{
		"src/resources/font/inter/inter.min.css",
		"src/resources/font/play/play.css",
		"src/resources/plugins/material-design-icons-iconfont/material-design-icons.min.css",
		"src/resources/plugins/fontawesome-free/css/all.min.css",
		"src/resources/plugins/simplebar/simplebar.min.css",
		"src/resources/plugins/summernote/summernote-bs4.css",
		"src/resources/plugins/select2/css/select2.min.css",
		"src/resources/plugins/flatpickr/flatpickr.min.css",
		"src/resources/plugins/noty/noty.css",
		"src/resources/plugins/noty/themes/relax.css",
		"src/resources/plugins/fancybox/fancybox.css",
		"src/resources/admin/css/style.css",
		"src/resources/admin/css/sidebar-dark.min.css",
		"src/resources/plugins/sweetalert2/sweetalert2.min.css",
	}
	adminJS := []string{
		"src/resources/plugins/jquery-3-6-0/jquery.min.js",
		"src/resources/plugins/bootstrap-4-6-1/js/bootstrap.bundle.js",
		"src/resources/plugins/simplebar/simplebar.min.js",
		"src/resources/plugins/feather-icons/feather.min.js",
		"src/resources/plugins/summernote/summernote-bs4.min.js",
		"src/resources/plugins/select2/js/select2.full.js",
		"src/resources/plugins/select2/js/i18n/ru.js",
		"src/resources/plugins/flatpickr/flatpickr.js",
		"src/resources/plugins/flatpickr/l10n/ru.js",
		"src/resources/plugins/noty/noty.js",
		"src/resources/plugins/inputmask/jquery.inputmask.js",
		"src/resources/plugins/fancybox/fancybox.umd.js",
		"src/resources/plugins/sweetalert2/sweetalert2.all.min.js",
		"src/resources/plugins/sortablejs/Sortable.min.js",
		"src/resources/plugins/js/script.min.js",
		"src/resources/plugins/date-format/date-format.js",
		"src/resources/plugins/npm/codemirror-bundle.js",
	}

	CSS := []string{
		"src/resources/plugins/bootstrap-5.0.2-dist/css/bootstrap.min.css",
	}
	JS := []string{
		"src/resources/plugins/jquery-3-6-0/jquery.min.js",
		"src/resources/plugins/bootstrap-5.0.2-dist/js/bootstrap.min.js",
		"src/public/admin/glob.js",
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
	mergeAndMinifyFiles(m, "text/css", adminCSS, "src/public/admin/app.css")
	mergeAndMinifyFiles(m, "application/javascript", adminJS, "src/public/admin/app.js")

	mergeAndMinifyFiles(m, "text/css", CSS, "src/public/app.css")
	mergeAndMinifyFiles(m, "application/javascript", JS, "src/public/app.js")
}

func minifyFile(m *minify.M, mediaType, inputPath, outputPath string) {
	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Ошибка открытия файла %s: %v", inputPath, err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Ошибка создания файла %s: %v", outputPath, err)
	}
	defer output.Close()

	if err := m.Minify(mediaType, output, input); err != nil {
		log.Fatalf("Ошибка минификации файла %s: %v", inputPath, err)
	}
}

func mergeAndMinifyFiles(m *minify.M, mediaType string, inputPaths []string, outputPath string) {
	var buffer bytes.Buffer

	for _, inputPath := range inputPaths {
		input, err := ioutil.ReadFile(inputPath)
		if err != nil {
			log.Fatalf("Ошибка чтения файла %s: %v", inputPath, err)
		}
		buffer.Write(input)
		buffer.WriteString("\n")
	}

	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Ошибка создания файла %s: %v", outputPath, err)
	}
	defer output.Close()

	if err := m.Minify(mediaType, output, &buffer); err != nil {
		log.Fatalf("Ошибка минификации файла %s: %v", outputPath, err)
	}
}
