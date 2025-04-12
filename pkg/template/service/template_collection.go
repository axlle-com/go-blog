package service

import (
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/template/models"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
	"sync"
)

type TemplateCollectionService struct {
	templateService *TemplateService
	templateRepo    templateRepository.TemplateRepository
	userProvider    user.UserProvider
}

func NewInfoBlockCollectionService(
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

func (s *TemplateCollectionService) WithPaginate(p contracts2.Paginator, filter *models.TemplateFilter) ([]*models.Template, error) {
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

	var users map[uint]contracts2.User

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

func (s *TemplateCollectionService) AggregatesResponses(infoBlocks []*InfoBlockResponse) []*InfoBlockResponse {
	var templateIDs []uint
	var userIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
	infoBlocksInterface := make([]contracts2.Resource, 0, len(infoBlocks))

	for _, infoBlock := range infoBlocks {
		infoBlocksInterface = append(infoBlocksInterface, infoBlock)
		if infoBlock.TemplateID != nil {
			id := *infoBlock.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}
		if infoBlock.UserID != nil {
			id := *infoBlock.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contracts2.User
	var templates map[uint]contracts2.Template
	var galleries map[uuid.UUID][]contracts2.Gallery

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

	for _, infoBlock := range infoBlocks {
		if gallery, ok := galleries[infoBlock.UUID]; ok {
			infoBlock.Galleries = gallery
		}
		if infoBlock.TemplateID != nil {
			infoBlock.Template = templates[*infoBlock.TemplateID]
		}
		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return infoBlocks
}
