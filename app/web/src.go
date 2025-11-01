package web

import (
	"bytes"
	"log"
	"os"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

func Minify(config contract.Config) {
	if !config.IsLocal() {
		logger.Info("[web][Minify] running on non-local environment; skipping HTML build")
		return
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
	minifyAdmin(m)
	if config.Layout() == "spring" {
		minifySpring(m)
	} else {
		minifyFront(m)
	}
}

func minifyAdmin(minify *minify.M) {
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

	mergeAndMinifyFiles(minify, "text/css", adminCSS, "src/public/admin/app.css")
	mergeAndMinifyFiles(minify, "application/javascript", adminJS, "src/public/admin/app.js")
}

func minifyFront(minify *minify.M) {
	CSS := []string{
		"src/resources/plugins/bootstrap-5.0.2-dist/css/bootstrap.min.css",
	}
	JS := []string{
		"src/resources/plugins/jquery-3-6-0/jquery.min.js",
		"src/resources/plugins/bootstrap-5.0.2-dist/js/bootstrap.min.js",
		"src/public/admin/glob.js",
	}

	mergeAndMinifyFiles(minify, "text/css", CSS, "src/public/app.css")
	mergeAndMinifyFiles(minify, "application/javascript", JS, "src/public/app.js")
}

func minifySpring(minify *minify.M) {
	CSS := []string{
		"src/resources/spring/css/bootstrap.css",
		"src/resources/spring/css/font-awesome.css",
		"src/resources/spring/css/themify-icons.css",
		"src/resources/spring/css/linear-icons.css",
		"src/resources/spring/css/animate.css",
		"src/resources/spring/css/owl.css",
		"src/resources/spring/css/jquery.fancybox.css",
		"src/resources/spring/css/responsive.css",
		"src/resources/spring/css/style.css",
	}
	JS := []string{
		"src/resources/spring/js/jquery.js",
		"src/resources/spring/js/bootstrap.min.js",
		"src/resources/spring/js/pagenav.js",
		"src/resources/spring/js/jquery.scrollTo.js",
		"src/resources/spring/js/jquery.fancybox.pack.js",
		"src/resources/spring/js/owl.js",
		"src/resources/spring/js/wow.js",
		"src/resources/spring/js/validate.js",
		"src/resources/spring/js/script.js",
		"src/public/admin/glob.js",
	}

	mergeAndMinifyFiles(minify, "text/css", CSS, "src/public/spring/css/app.css")
	mergeAndMinifyFiles(minify, "application/javascript", JS, "src/public/spring/js/app.js")
}

func mergeAndMinifyFiles(minifyTool *minify.M, mediaType string, inputPaths []string, outputPath string) {
	var buffer bytes.Buffer

	for _, inputPath := range inputPaths {
		input, err := os.ReadFile(inputPath)
		if err != nil {
			logger.Fatalf("[web][mergeAndMinifyFiles] error reading file %s: %v", inputPath, err)
		}
		buffer.Write(input)
		buffer.WriteString("\n")
	}

	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("[web][mergeAndMinifyFiles] error creating file %s: %v", outputPath, err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			logger.Errorf("[web][mergeAndMinifyFiles] close error: %v", err)
		}
	}(output)

	if err := minifyTool.Minify(mediaType, output, &buffer); err != nil {
		logger.Fatalf("[web][mergeAndMinifyFiles] minify error for file %s: %v", outputPath, err)
	}
}

func minifyFile(m *minify.M, mediaType, inputPath, outputPath string) {
	input, err := os.Open(inputPath)
	if err != nil {
		logger.Fatalf("[web][minifyFile] open error for file %s: %v", inputPath, err)
	}
	defer func(input *os.File) {
		err := input.Close()
		if err != nil {
			logger.Errorf("[web][minifyFile] close error: %v", err)
		}
	}(input)

	output, err := os.Create(outputPath)
	if err != nil {
		logger.Fatalf("[web][minifyFile] create error for file %s: %v", outputPath, err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			logger.Errorf("[web][minifyFile] close error: %v", err)
		}
	}(output)

	if err := m.Minify(mediaType, output, input); err != nil {
		logger.Fatalf("[web][minifyFile] minify error for file %s: %v", inputPath, err)
	}
}
