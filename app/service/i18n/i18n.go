package i18n

import (
	"encoding/json"
	"path/filepath"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Service struct {
	bundle *i18n.Bundle
}

func New(cfg contract.Config, diskService contract.DiskService) *Service {
	newBundle := i18n.NewBundle(language.English)
	newBundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Загружаем файлы через DiskService (из embed или с диска)
	for _, lang := range []string{"en", "ru"} {
		path := filepath.Join("i18n", "locales", lang+".json")

		// Читаем файл через DiskService
		data, err := diskService.ReadFile(path)
		if err != nil {
			logger.Errorf("[I18n] Failed to read locale file %s: %v", path, err)
			panic(err)
		}

		// Парсим и загружаем сообщения из байтов
		_, err = newBundle.ParseMessageFileBytes(data, path)
		if err != nil {
			logger.Errorf("[I18n] Failed to parse locale file %s: %v", path, err)
			panic(err)
		}
		logger.Infof("[I18n] Loaded locale file: %s", path)
	}
	return &Service{bundle: newBundle}
}

func (s *Service) Localizer(langs ...string) *i18n.Localizer {
	var tags []string

	for _, lang := range langs {
		if lang != "" {
			tags = append(tags, lang)
		}
	}
	tags = []string{"en"}

	return i18n.NewLocalizer(s.bundle, tags...)
}
