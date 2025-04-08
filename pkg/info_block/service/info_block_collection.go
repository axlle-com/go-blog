package service

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
	"sync"
)

type InfoBlockCollectionService struct {
	infoBlockRepo    repository.InfoBlockRepository
	resourceRepo     repository.InfoBlockHasResourceRepository
	galleryProvider  provider.GalleryProvider
	templateProvider template.TemplateProvider
	userProvider     user.UserProvider
}

func NewInfoBlockCollectionService(
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
	galleryProvider provider.GalleryProvider,
	templateProvider template.TemplateProvider,
	userProvider user.UserProvider,
) *InfoBlockCollectionService {
	return &InfoBlockCollectionService{
		infoBlockRepo:    infoBlockRepo,
		resourceRepo:     resourceRepo,
		galleryProvider:  galleryProvider,
		templateProvider: templateProvider,
		userProvider:     userProvider,
	}
}

func (s *InfoBlockCollectionService) GetAll() ([]*InfoBlock, error) {
	return s.infoBlockRepo.GetAll()
}

func (s *InfoBlockCollectionService) DeleteInfoBlocks(infoBlocks []*InfoBlock) (err error) {
	var ids []uint
	for _, infoBlock := range infoBlocks {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.infoBlockRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}

func (s *InfoBlockCollectionService) WithPaginate(p contracts.Paginator, filter *InfoBlockFilter) ([]*InfoBlock, error) {
	return s.infoBlockRepo.WithPaginate(p, filter)
}

func (s *InfoBlockCollectionService) Aggregates(infoBlocks []*InfoBlock) []*InfoBlock {
	var templateIDs []uint
	var userIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)

	for _, infoBlock := range infoBlocks {
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

	var users map[uint]contracts.User
	var templates map[uint]contracts.Template

	wg.Add(2)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.templateProvider.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

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
		if infoBlock.TemplateID != nil {
			infoBlock.Template = templates[*infoBlock.TemplateID]
		}
		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return infoBlocks
}

func (s *InfoBlockCollectionService) AggregatesResponses(infoBlocks []*InfoBlockResponse) []*InfoBlockResponse {
	var templateIDs []uint
	var userIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
	infoBlocksInterface := make([]contracts.Resource, 0, len(infoBlocks))

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

	var users map[uint]contracts.User
	var templates map[uint]contracts.Template
	var galleries map[uuid.UUID][]contracts.Gallery

	wg.Add(3)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.templateProvider.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

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

	go func() {
		defer wg.Done()
		galleries = s.galleryProvider.GetIndexesForResources(infoBlocksInterface)
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
