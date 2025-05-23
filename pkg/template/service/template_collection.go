package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/template/models"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"sync"
)

type TemplateCollectionService struct {
	templateService *TemplateService
	templateRepo    templateRepository.TemplateRepository
	userProvider    user.UserProvider
}

func NewTemplateCollectionService(
	templateService *TemplateService,
	templateRepo templateRepository.TemplateRepository,
	userProvider user.UserProvider,
) *TemplateCollectionService {
	return &TemplateCollectionService{
		templateService: templateService,
		templateRepo:    templateRepo,
		userProvider:    userProvider,
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

func (s *TemplateCollectionService) WithPaginate(p contracts.Paginator, filter *models.TemplateFilter) ([]*models.Template, error) {
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

	var users map[uint]contracts.User

	wg.Add(1)

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.userProvider.GetMapByIDs(userIDs)
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
