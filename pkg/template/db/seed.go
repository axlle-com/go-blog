package db

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	template repository.TemplateRepository
}

func NewSeeder(
	template repository.TemplateRepository,
) contract.Seeder {
	return &seeder{
		template: template,
	}
}

func (s *seeder) Seed() error {
	cfg := config.Config()

	// templates/<layout>
	layout := cfg.Layout()
	templatesRoot := cfg.SrcFolderBuilder(filepath.Join("templates", "front", layout))

	// Если папки нет — ничего не делаем (не падаем)
	if _, err := os.Stat(templatesRoot); err != nil {
		if os.IsNotExist(err) {
			logger.Infof("[template][seeder][Seed] templates root not found: %s", templatesRoot)
			return nil
		}
		return err
	}

	// Обходим файлы темы и добавляем все .gohtml кроме default.gohtml
	walkErr := filepath.Walk(templatesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".gohtml" {
			return nil
		}

		base := filepath.Base(path)
		if strings.EqualFold(base, "default.gohtml") {
			return nil
		}

		// Имя шаблона = имя файла без расширения
		name := strings.TrimSuffix(base, filepath.Ext(base))

		// Ресурс — это директория сразу под layout, например: templates/<layout>/<resource>/<file>.gohtml
		rel, err := filepath.Rel(templatesRoot, path)
		if err != nil {
			return err
		}
		// Получаем имя корневой папки ресурса
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) < 2 {
			// Файлы непосредственно в корне layout пропускаем
			return nil
		}
		resourceName := parts[0]

		// Читаем содержимое шаблона и сохраняем в HTML
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		html := string(data)

		// Пытаемся извлечь заголовок из первой строки-комментария
		title := extractTitleFromTemplate(html)
		if title == "" {
			title = fmt.Sprintf("%s/%s", resourceName, name)
		}

		// Сначала проверяем существование шаблона по theme, name и resource_name через фильтр
		filter := models.NewTemplateFilter()
		filter.Theme = &layout
		filter.Name = &name
		filter.ResourceName = &resourceName
		existing, err := s.template.FindByFilter(filter)
		if err == nil && existing != nil {
			// Шаблон уже существует, пропускаем без ошибки
			logger.Infof("[template][seeder][Seed] template with theme=%s, name=%s and resource_name=%s already exists (ID=%d), skipping", layout, name, resourceName, existing.ID)
			return nil
		}

		now := time.Now()
		tpl := models.Template{
			Title:        title,
			IsMain:       false,
			Name:         name,
			Theme:        db.StrPtr(layout),
			ResourceName: db.StrPtr(resourceName),
			HTML:         db.StrPtr(html),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		// Создаем шаблон, если он не найден
		if err := s.template.Create(&tpl); err != nil {
			logger.Errorf("[template][seeder][Seed] error creating template %s: %v", tpl.Name, err)
			return err
		}
		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	logger.Info("[template][seeder][Seed] Database seeded Template from filesystem successfully!")
	return nil
}

func (s *seeder) SeedTest(n int) error {
	resources := map[string]string{
		"Шаблон для постов":    "posts",
		"Шаблон для категорий": "post_categories",
		"Шаблон для тегов":     "post_tags",
		"Шаблон для блоков":    "info_blocks",
		"Шаблон для меню":      "menus",
	}

	// один раз посеять rng
	rand.Seed(time.Now().UnixNano())

	// подготовим слайс ключей
	keys := make([]string, 0, len(resources))
	for title := range resources {
		keys = append(keys, title)
	}

	for i := 1; i <= n; i++ {
		var template models.Template
		now := time.Now()

		// берём случайный ключ и значение из map
		title := keys[rand.Intn(len(keys))]
		resource := resources[title]

		template.Title = fmt.Sprintf("%s #%d", title, i)
		template.Name = faker.Username()
		template.Theme = db.StrPtr(config.Config().Layout())
		template.ResourceName = db.StrPtr(resource)
		template.CreatedAt = &now
		template.UpdatedAt = &now

		if err := s.template.Create(&template); err != nil {
			return err
		}
	}

	logger.Info("[template][seeder][SeedTest] Database seeded Template successfully!")
	return nil
}

// extractTitleFromTemplate возвращает текст из первой непустой строки-комментария
// вида <!-- ... -->. Если такой строки нет — возвращает пустую строку.
func extractTitleFromTemplate(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "<!--") && strings.Contains(trimmed, "-->") {
			inner := strings.TrimPrefix(trimmed, "<!--")
			inner = strings.TrimSuffix(inner, "-->")
			return strings.TrimSpace(inner)
		}
		// первая непустая строка не комментарий — выходим
		break
	}
	return ""
}
