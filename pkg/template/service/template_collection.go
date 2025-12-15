package service

import (
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
)

type TemplateCollectionService struct {
	templateService *TemplateService
	templateRepo    repository.TemplateRepository
	api             *api.Api
}

func NewTemplateCollectionService(
	templateService *TemplateService,
	templateRepo repository.TemplateRepository,
	api *api.Api,
) *TemplateCollectionService {
	return &TemplateCollectionService{
		templateService: templateService,
		templateRepo:    templateRepo,
		api:             api,
	}
}

func (s *TemplateCollectionService) GetAll() ([]*models.Template, error) {
	return s.templateRepo.GetAll()
}

func (s *TemplateCollectionService) DeleteTemplates(templates []*models.Template) (err error) {
	var ids []uint
	for _, infoBlock := range templates {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.templateRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}

func (s *TemplateCollectionService) WithPaginate(p contract.Paginator, filter *request.TemplateFilter) ([]*models.Template, error) {
	return s.templateRepo.WithPaginate(p, filter)
}

func (s *TemplateCollectionService) Aggregates(templates []*models.Template) []*models.Template {
	var userIDs []uint

	userIDsMap := make(map[uint]bool)

	for _, template := range templates {
		if template.UserID != nil {
			id := *template.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User

	wg.Add(1)

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, infoBlock := range templates {
		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return templates
}
