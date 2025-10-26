package service

import (
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/template/provider"
	"github.com/google/uuid"
)

type TagCollectionService struct {
	tagService       *TagService
	tagRepo          repository.PostTagRepository
	resourceRepo     repository.PostTagResourceRepository
	templateProvider provider.TemplateProvider
}

func NewTagCollectionService(
	tagService *TagService,
	postTagRepo repository.PostTagRepository,
	resourceRepo repository.PostTagResourceRepository,
	templateProvider provider.TemplateProvider,
) *TagCollectionService {
	return &TagCollectionService{
		tagService:       tagService,
		tagRepo:          postTagRepo,
		resourceRepo:     resourceRepo,
		templateProvider: templateProvider,
	}
}

func (s *TagCollectionService) WithPaginate(p contracts.Paginator, filter *models.TagFilter) ([]*models.PostTag, error) {
	return s.tagRepo.WithPaginate(p, filter)
}

func (s *TagCollectionService) Aggregates(tags []*models.PostTag) []*models.PostTag {
	var templateIDs []uint

	templateIDsMap := make(map[uint]bool)

	for _, tag := range tags {
		if tag.TemplateID != nil {
			id := *tag.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}
	}

	var templates map[uint]contracts.Template

	if len(templateIDs) > 0 {
		var err error
		templates, err = s.templateProvider.GetMapByIDs(templateIDs)
		if err != nil {
			logger.Error(err)
		}
	}

	for _, tag := range tags {
		if tag.TemplateID != nil {
			tag.Template = templates[*tag.TemplateID]
		}
	}

	return tags
}

func (s *TagCollectionService) SyncTags(resourceUUID uuid.UUID, namees []string) ([]*models.PostTag, error) {
	// 0. Нормализуем входящие заголовки
	var allNames []string
	for _, title := range namees {
		t := strings.TrimSpace(title)
		if t != "" {
			allNames = append(allNames, t)
		}
	}

	// 1. Загрузить существующие теги по Name
	tagMap := make(map[string]*models.PostTag, len(allNames))
	found, err := s.tagRepo.GetByNames(allNames)
	if err != nil {
		return nil, err
	}
	for _, t := range found {
		tagMap[t.Name] = t
	}

	// 2. Сформировать полный список тегов (существующих + новых)
	var allTags []*models.PostTag
	for _, name := range allNames {
		if existing, ok := tagMap[name]; ok {
			allTags = append(allTags, existing)
		} else {
			raw := &models.PostTag{Name: name}
			created, err := s.tagService.Create(raw)
			if err != nil {
				return nil, err
			}
			allTags = append(allTags, created)
		}
	}

	// 3. Получить текущие связи ресурс–тег
	currentRelations, err := s.resourceRepo.GetForResource(resourceUUID)
	if err != nil {
		return nil, err
	}

	currentIDs := make(map[uint]struct{}, len(currentRelations))
	for _, rel := range currentRelations {
		currentIDs[rel.PostTagID] = struct{}{}
	}

	// 4. Вычислить, что отвязать, а что привязать
	wantIDs := make(map[uint]struct{}, len(allTags))
	for _, tag := range allTags {
		wantIDs[tag.ID] = struct{}{}
	}

	// 4a. Отвязать лишние
	for id := range currentIDs {
		if _, needed := wantIDs[id]; !needed {
			if err := s.resourceRepo.DeleteByParams(resourceUUID, id); err != nil {
				return nil, err
			}
		}
	}

	// 4b. Привязать недостающие
	for id := range wantIDs {
		if _, exists := currentIDs[id]; !exists {
			rel := &models.PostTagHasResource{
				ResourceUUID: resourceUUID,
				PostTagID:    id,
			}
			if err := s.resourceRepo.Create(rel); err != nil {
				return nil, err
			}
		}
	}

	// 5. Возвращаем актуальный слайс тегов для ресурса
	return allTags, nil
}

func (s *TagCollectionService) GetAll() ([]*models.PostTag, error) {
	return s.tagRepo.GetAll()
}

func (s *TagCollectionService) GetForResource(resource contracts.Resource) ([]*models.PostTag, error) {
	return s.tagRepo.GetForResource(resource)
}

func (s *TagCollectionService) UpdateInfoBlockSnapshots(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	return s.tagRepo.UpdateFieldsByUUIDs(uuids, patch)
}
