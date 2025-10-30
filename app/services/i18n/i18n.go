package i18n

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Service struct {
	bundle *i18n.Bundle
}

func New(cfg contract.Config) *Service {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Используем конфиг для получения пути к файлам локализации
	localesDir := cfg.SrcFolderBuilder("services/i18n/locales")

	// Загружаем файлы из файловой системы
	for _, lang := range []string{"en", "ru"} {
		path := filepath.Join(localesDir, lang+".json")
		if _, err := b.LoadMessageFile(path); err != nil {
			logger.Errorf("[I18n] Failed to load locale file %s: %v", path, err)
			panic(err)
		}
		logger.Infof("[I18n] Loaded locale file: %s", path)
	}
	return &Service{bundle: b}
}

// Localizer строит локализатор по приоритету: query?lang → cookie(lang) → Accept-Language → en.
func (s *Service) Localizer(r *http.Request) *i18n.Localizer {
	accept := r.Header.Get("Accept-Language")
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		if c, err := r.Cookie("lang"); err == nil {
			lang = c.Value
		}
	}

	var tags []string
	if lang != "" {
		tags = append(tags, lang)
	}
	if accept != "" {
		tags = append(tags, accept)
	}
	tags = append(tags, "en")

	return i18n.NewLocalizer(s.bundle, tags...)
}
