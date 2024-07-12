package web

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"io/ioutil"
	"log"
	"os"
)

func InitMinify() {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	mergeAndMinifyFiles(m, "text/css", []string{"src/css/bootstrap.min.css"}, "src/app.css")
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
