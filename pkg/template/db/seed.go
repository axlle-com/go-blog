package db

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path"
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
	layout := strings.TrimSpace(s.config.Layout())
	if layout == "" {
		layout = "default"
	}

	// templates/front/<layout>
	templatesRoot := path.Join("templates", "front", layout)

	templatesFS := s.disk.GetTemplatesFS()

	// Если папки нет — ничего не делаем (не падаем)
	if _, err := fs.Stat(templatesFS, templatesRoot); err != nil {
		if isNotExistFS(err) {
			logger.Errorf("[template][seeder][Seed] templates root not found in FS: %s", templatesRoot)
			return nil
		}

		return fmt.Errorf("stat templates root %q: %w", templatesRoot, err)
	}

	walkErr := fs.WalkDir(templatesFS, templatesRoot, func(pathString string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(path.Ext(pathString))
		if ext != ".gohtml" {
			return nil
		}

		base := path.Base(pathString)
		if strings.EqualFold(base, "default.gohtml") {
			return nil
		}

		// Имя шаблона = имя файла без расширения
		name := strings.TrimSuffix(base, ext)

		// Ресурс — это директория сразу под layout:
		// templates/front/<layout>/<resource>/<file>.gohtml
		rel := strings.TrimPrefix(pathString, templatesRoot+"/")
		parts := strings.Split(rel, "/")
		if len(parts) < 2 {
			// Файлы непосредственно в корне layout пропускаем
			return nil
		}
		resourceName := parts[0]

		// Читаем содержимое шаблона из templatesFS
		data, err := fs.ReadFile(templatesFS, pathString)
		if err != nil {
			return err
		}
		html := string(data)

		// Пытаемся извлечь заголовок из первой строки-комментария
		title := extractTitleFromTemplate(html)
		if title == "" {
			title = fmt.Sprintf("%s/%s", resourceName, name)
		}

		// Проверяем существование по фильтру (важно: указатели на копии)
		filter := models.NewTemplateFilter()
		layoutCopy := layout
		nameCopy := name
		resourceCopy := resourceName

		filter.Theme = &layoutCopy
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

	layout := s.config.Layout()

	for i := 1; i <= n; i++ {
		now := time.Now()

		title := keys[rand.Intn(len(keys))]
		resource := resources[title]

		tpl := models.Template{
			Title:        fmt.Sprintf("%s #%d", title, i),
			Name:         faker.Username(),
			Theme:        db.StrPtr(layout),
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
