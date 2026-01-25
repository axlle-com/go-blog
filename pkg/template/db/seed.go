package db

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	config   contract.Config
	disk     contract.DiskService
	template repository.TemplateRepository
}

func NewSeeder(
	config contract.Config,
	disk contract.DiskService,
	template repository.TemplateRepository,
) contract.Seeder {
	return &seeder{
		config:   config,
		disk:     disk,
		template: template,
	}
}

func (s *seeder) Seed() error {
	layouts, err := s.listFrontLayouts()
	if err != nil {
		return err
	}
	if len(layouts) == 0 {
		logger.Infof("[template][seeder][Seed] no layouts found in templates/front")
		return nil
	}

	for _, layout := range layouts {
		if err := s.seedLayout(layout); err != nil {
			return err
		}
	}

	logger.Info("[template][seeder][Seed] Database seeded Template from FS successfully!")

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

	rand.Seed(time.Now().UnixNano())

	keys := make([]string, 0, len(resources))
	for title := range resources {
		keys = append(keys, title)
	}

	for i := 1; i <= n; i++ {
		now := time.Now()

		title := keys[rand.Intn(len(keys))]
		resource := resources[title]

		tpl := models.Template{
			Title:        fmt.Sprintf("%s #%d", title, i),
			Name:         faker.Username(),
			ResourceName: db.StrPtr(resource),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		if err := s.template.Create(&tpl); err != nil {
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

// extractDefineNameFromTemplate возвращает имя из первого {{ define "..." }} в шаблоне.
func extractDefineNameFromTemplate(content string) string {
	re := regexp.MustCompile(`(?i)\{\{\s*define\s+"([^"]+)"\s*\}\}`)
	match := re.FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}

	return strings.TrimSpace(match[1])
}

func isNotExistFS(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, fs.ErrNotExist) {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "file does not exist") || strings.Contains(msg, "no such file")
}

func (s *seeder) listFrontLayouts() ([]string, error) {
	frontRoot := path.Join("templates", "front")

	entries, err := s.disk.ReadDir(frontRoot)
	if err != nil {
		if isNotExistFS(err) {
			logger.Infof("[template][seeder][Seed] templates root not found: %s", frontRoot)
			return nil, nil
		}

		return nil, fmt.Errorf("read templates root %q: %w", frontRoot, err)
	}

	layouts := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			layouts = append(layouts, entry.Name())
		}
	}

	return layouts, nil
}

func (s *seeder) seedLayout(layout string) error {
	templatesRoot := path.Join("templates", "front", layout)

	walkErr := s.walkTemplates(templatesRoot, func(pathString string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(path.Ext(pathString))
		if ext != ".gohtml" {
			return nil
		}

		if strings.EqualFold(path.Base(pathString), "default.gohtml") {
			return nil
		}

		// Ресурс — это директория сразу под layout:
		// templates/front/<layout>/<resource>/<path>/<file>.gohtml
		rel := strings.TrimPrefix(pathString, templatesRoot+"/")
		relNoExt := strings.TrimSuffix(rel, ext)
		parts := strings.Split(relNoExt, "/")
		if len(parts) < 2 {
			// Файлы непосредственно в корне layout пропускаем
			return nil
		}

		resourceName := parts[0]
		name := fmt.Sprintf("%s.%s", layout, strings.ReplaceAll(relNoExt, "/", "."))

		data, err := s.disk.ReadFile(pathString)
		if err != nil {
			return err
		}

		html := string(data)
		if defineName := extractDefineNameFromTemplate(html); defineName != "" {
			name = defineName
		}

		// Пытаемся извлечь заголовок из первой строки-комментария
		title := extractTitleFromTemplate(html)
		if title == "" {
			title = name
		}

		// Проверяем существование по фильтру (важно: указатели на копии)
		filter := models.NewTemplateFilter()
		nameCopy := name
		resourceCopy := resourceName

		filter.Name = &nameCopy
		filter.ResourceName = &resourceCopy

		existing, err := s.template.FindByFilter(filter)
		if err == nil && existing != nil {
			logger.Infof(
				"[template][seeder][Seed] template theme=%s name=%s resource=%s already exists (ID=%d), skipping",
				layout, name, resourceName, existing.ID,
			)
			return nil
		}

		now := time.Now()
		tpl := models.Template{
			Title:        title,
			IsMain:       false,
			Name:         name,
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

	return nil
}

func (s *seeder) walkTemplates(root string, fn func(pathString string, info fs.DirEntry) error) error {
	entries, err := s.disk.ReadDir(root)
	if err != nil {
		if isNotExistFS(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		fullPath := path.Join(root, entry.Name())
		if entry.IsDir() {
			if err := s.walkTemplates(fullPath, fn); err != nil {
				return err
			}
			continue
		}

		if err := fn(fullPath, entry); err != nil {
			return err
		}
	}

	return nil
}
