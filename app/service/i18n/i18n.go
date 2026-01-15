package i18n

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Service struct {
	bundle    *i18n.Bundle
	supported map[string]struct{}
}

func New(cfg contract.Config, diskService contract.DiskService) *Service {
	_ = cfg

	newBundle := i18n.NewBundle(language.English)
	newBundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	svc := &Service{
		bundle: newBundle,
		supported: map[string]struct{}{
			"en": {},
			"ru": {},
		},
	}

	// Загружаем файлы через DiskService (из embed или с диска)
	for _, lang := range []string{"en", "ru"} {
		path := filepath.Join("services", "i18n", "locales", lang+".json")

		data, err := diskService.ReadFile(path)
		if err != nil {
			logger.Errorf("[I18n] Failed to read locale file %s: %v", path, err)
			panic(err)
		}

		if _, err := newBundle.ParseMessageFileBytes(data, path); err != nil {
			logger.Errorf("[I18n] Failed to parse locale file %s: %v", path, err)
			panic(err)
		}

		logger.Infof("[I18n] Loaded locale file: %s", path)
	}

	return svc
}

// NormalizeLang принимает "ru", "ru-RU", "ru_RU", "en-US" и т.п.
// Возвращает нормализованный BCP-47 tag (например "ru", "en-US") и ok=true,
// только если базовый язык поддерживается сервисом.
func (s *Service) NormalizeLang(in string) (tag string, ok bool) {
	in = strings.TrimSpace(in)
	if in == "" {
		return "", false
	}

	// ru_RU -> ru-RU
	in = strings.ReplaceAll(in, "_", "-")

	// Если передали Accept-Language целиком, например:
	// "ru-RU,ru;q=0.9,en;q=0.8" — берём только первый токен.
	in = strings.TrimSpace(strings.Split(in, ",")[0])

	t, err := language.Parse(in)
	if err != nil {
		return "", false
	}

	base, _ := t.Base()
	baseStr := strings.ToLower(base.String())
	if _, exists := s.supported[baseStr]; !exists {
		return "", false
	}

	// Каноничный вид тега
	return t.String(), true
}

// Localizer собирает приоритеты языков из аргументов.
// Можно передавать сюда и "lang" и "Accept-Language" — сервис сам отфильтрует.
func (s *Service) Localizer(langs ...string) *i18n.Localizer {
	tags := make([]string, 0, len(langs)+2)
	seen := make(map[string]struct{}, 4)

	add := func(x string) {
		if x == "" {
			return
		}
		if _, ok := seen[x]; ok {
			return
		}
		seen[x] = struct{}{}
		tags = append(tags, x)
	}

	for _, l := range langs {
		tag, ok := s.NormalizeLang(l)
		if !ok {
			continue
		}

		add(tag)
		if base := strings.ToLower(strings.SplitN(tag, "-", 2)[0]); base != "" {
			add(base)
		}
	}

	add("en")

	return i18n.NewLocalizer(s.bundle, tags...)
}
