package web

import (
	"bytes"
	"fmt"
	"os"

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
	if !w.cfg.IsLocal() {
		logger.Info("[web][Minify] running on non-local environment; skipping HTML build")
		return nil
	}

	if err := w.minifyAdmin(); err != nil {
		return fmt.Errorf("minify admin: %w", err)
	}

	if w.cfg.Layout() == "spring" {
		if err := w.minifySpring(); err != nil {
			return fmt.Errorf("minify spring: %w", err)
		}
	} else {
		if err := w.minifyFront(); err != nil {
			return fmt.Errorf("minify front: %w", err)
		}
	}

	return nil
}

func (w *WebMinifier) minifyAdmin() error {
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
	}

	if err := w.mergeAndMinifyFiles("text/css", adminCSS, "src/public/admin/app.css"); err != nil {
		return err
	}

	if err := w.mergeAndMinifyFiles("application/javascript", adminJS, "src/public/admin/app.js"); err != nil {
		return err
	}

	return nil
}

func (w *WebMinifier) minifyFront() error {
	CSS := []string{
		"src/resources/plugins/bootstrap-5.0.2-dist/css/bootstrap.min.css",
	}
	JS := []string{
		"src/resources/plugins/jquery-3-6-0/jquery.min.js",
		"src/resources/plugins/bootstrap-5.0.2-dist/js/bootstrap.min.js",
		"src/public/admin/glob.js",
	}

	if err := w.mergeAndMinifyFiles("text/css", CSS, "src/public/app.css"); err != nil {
		return err
	}
	if err := w.mergeAndMinifyFiles("application/javascript", JS, "src/public/app.js"); err != nil {
		return err
	}
	return nil
}

func (w *WebMinifier) minifySpring() error {
	CSS := []string{
		"src/public/spring/css/font.css",
		"src/resources/spring/css/bootstrap.css",
		"src/resources/spring/css/font-awesome.css",
		"src/resources/spring/css/themify-icons.css",
		"src/resources/spring/css/linear-icons.css",
		"src/resources/spring/css/animate.css",
		"src/resources/spring/css/owl.css",
		"src/resources/spring/css/jquery.fancybox.css",
		"src/resources/spring/css/style.css",
		"src/resources/spring/css/responsive.css",
		"src/public/spring/css/common.css",
	}
	JS := []string{
		"src/resources/spring/js/jquery.js",
		"src/resources/spring/js/bootstrap.min.js",
		"src/resources/spring/js/appear.js",
		"src/resources/spring/js/pagenav.js",
		"src/resources/spring/js/jquery.scrollTo.js",
		"src/resources/spring/js/jquery.fancybox.pack.js",
		"src/resources/spring/js/owl.js",
		"src/resources/spring/js/wow.js",
		"src/resources/spring/js/validate.js",
		"src/resources/spring/js/script.js",
		"src/public/admin/glob.js",
	}

	if err := w.mergeAndMinifyFiles("text/css", CSS, "src/public/spring/css/app.css"); err != nil {
		return err
	}

	if err := w.mergeAndMinifyFiles("application/javascript", JS, "src/public/spring/js/app.js"); err != nil {
		return err
	}

	return nil
}

func (w *WebMinifier) mergeAndMinifyFiles(mediaType string, inputPaths []string, outputPath string) error {
	var buffer bytes.Buffer

	for _, inputPath := range inputPaths {
		input, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("[mergeAndMinifyFiles] read %s: %w", inputPath, err)
		}

		buffer.Write(input)
		buffer.WriteString("\n")
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

	if err := w.m.Minify(mediaType, output, &buffer); err != nil {
		return fmt.Errorf("[mergeAndMinifyFiles] minify %s: %w", outputPath, err)
	}

	return nil
}
