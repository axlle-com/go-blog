package service

import (
	"sort"
	"strings"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
)

type InfoBlockCollectionService struct {
	infoBlockRepo repository.InfoBlockRepository
	resourceRepo  repository.InfoBlockHasResourceRepository
	api           *api.Api
}

func NewInfoBlockCollectionService(
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
	api *api.Api,
) *InfoBlockCollectionService {
	return &InfoBlockCollectionService{
		infoBlockRepo: infoBlockRepo,
		resourceRepo:  resourceRepo,
		api:           api,
	}
}

func (s *InfoBlockCollectionService) GetAll() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAll()
}

func (s *InfoBlockCollectionService) GetRoots() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetRoots()
}

func (s *InfoBlockCollectionService) GetForResourceByFilter(filter *models.InfoBlockFilter) []*models.InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockCollectionService][GetForResourceByFilter] Error: %v", err)
		return nil
	}

	return s.AggregatesResponses(infoBlocks)
}

func (s *InfoBlockCollectionService) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAllForParent(parent)
}

func (s *InfoBlockCollectionService) WithPaginate(paginator contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.WithPaginate(paginator, filter)
}

func (s *InfoBlockCollectionService) Aggregates(infoBlocks []*models.InfoBlock) []*models.InfoBlock {
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

	var users map[uint]contract.User
	var templates map[uint]contract.Template

	wg.Add(2)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Errorf("[info_block][InfoBlockCollectionService][Aggregates] Error: %v", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Errorf("[info_block][InfoBlockCollectionService][Aggregates] Error: %v", err)
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

func (s *InfoBlockCollectionService) AggregatesResponses(infoBlocks []*models.InfoBlockResponse) []*models.InfoBlockResponse {
	if len(infoBlocks) == 0 {
		return []*models.InfoBlockResponse{}
	}

	// 1) уникальные root IDs (сохраняя порядок входа)
	rootIDs := make([]uint, 0, len(infoBlocks))
	seenRoot := make(map[uint]struct{}, len(infoBlocks))
	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}
		if _, ok := seenRoot[ib.ID]; ok {
			continue
		}
		seenRoot[ib.ID] = struct{}{}
		rootIDs = append(rootIDs, ib.ID)
	}
	if len(rootIDs) == 0 {
		return []*models.InfoBlockResponse{}
	}

	// 2) пути корней: сначала из входа; если чего-то нет — добираем из БД
	paths := make([]string, 0, len(rootIDs))
	missingIDs := make([]uint, 0)

	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}
		p := strings.TrimSpace(ib.PathLtree)
		if p != "" {
			paths = append(paths, p)
		} else {
			missingIDs = append(missingIDs, ib.ID)
		}
	}

	if len(missingIDs) > 0 {
		rootModels, err := s.infoBlockRepo.GetByIDs(missingIDs)
		if err != nil {
			logger.Errorf("[info_block][InfoBlockCollectionService][AggregatesResponses] Error: %v", err)
			s.enrichInfoBlockResponses(infoBlocks)
			return infoBlocks
		}
		for _, m := range rootModels {
			if m != nil {
				p := strings.TrimSpace(m.PathLtree)
				if p != "" {
					paths = append(paths, p)
				}
			}
		}
	}

	paths = collapsePaths(paths)
	if len(paths) == 0 {
		s.enrichInfoBlockResponses(infoBlocks)
		return infoBlocks
	}

	// 3) одним запросом получаем все узлы поддеревьев (включая корни)
	allNodes, err := s.infoBlockRepo.GetSubtreesByPaths(paths)
	if err != nil {
		logger.Errorf("[info_block][InfoBlockCollectionService][AggregatesResponses] Error: %v", err)
		s.enrichInfoBlockResponses(infoBlocks)
		return infoBlocks
	}

	// 4) ID -> *InfoBlockResponse
	// Корни оставляем как есть (могут содержать relation_id/sort/position из GetForResourceByFilter)
	respByID := make(map[uint]*models.InfoBlockResponse, len(allNodes)+len(infoBlocks))
	for _, r := range infoBlocks {
		if r != nil {
			respByID[r.ID] = r
		}
	}

	// Создаём responses для остальных узлов поддерева
	for _, n := range allNodes {
		if n == nil {
			continue
		}

		if existing, ok := respByID[n.ID]; ok {
			// корень уже есть как response — докинем недостающее (например PathLtree)
			if existing != nil && existing.PathLtree == "" {
				existing.PathLtree = n.PathLtree
			}
			// Sort у корней лучше не трогать (может быть relation sort/position),
			// дети сортируются по своему Sort из InfoBlock.
			continue
		}

		respByID[n.ID] = &models.InfoBlockResponse{
			ID:          n.ID,
			UUID:        n.UUID,
			TemplateID:  n.TemplateID,
			InfoBlockID: n.InfoBlockID,
			UserID:      n.UserID,
			Media:       n.Media,
			Title:       n.Title,
			Description: n.Description,
			Image:       n.Image,
			PathLtree:   n.PathLtree,
			Sort:        n.Sort,
		}
	}

	// 5) parentID -> children[]
	childrenByParent := make(map[uint][]*models.InfoBlockResponse)
	for _, n := range allNodes {
		if n == nil || n.InfoBlockID == nil || *n.InfoBlockID == 0 {
			continue
		}
		parentID := *n.InfoBlockID
		childResp := respByID[n.ID]
		if childResp == nil {
			continue
		}
		childrenByParent[parentID] = append(childrenByParent[parentID], childResp)
	}

	// 6) сортируем детей каждого родителя по Sort ASC, потом ID ASC
	for pid := range childrenByParent {
		kids := childrenByParent[pid]
		sort.Slice(kids, func(i, j int) bool {
			if kids[i].Sort != kids[j].Sort {
				return kids[i].Sort < kids[j].Sort
			}
			return kids[i].ID < kids[j].ID
		})
		childrenByParent[pid] = kids
	}

	// 7) навешиваем Children рекурсивно
	var attach func(id uint) []*models.InfoBlockResponse
	attach = func(id uint) []*models.InfoBlockResponse {
		kids := childrenByParent[id]
		for _, ch := range kids {
			ch.Children = attach(ch.ID)
		}
		return kids
	}

	for _, root := range infoBlocks {
		if root == nil {
			continue
		}
		root.Children = attach(root.ID)
	}

	// 8) обогащаем ВСЕ ноды (корни + дети) одним пакетом
	allResponses := make([]*models.InfoBlockResponse, 0, len(respByID))
	for _, r := range respByID {
		if r != nil {
			allResponses = append(allResponses, r)
		}
	}
	s.enrichInfoBlockResponses(allResponses)

	return infoBlocks
}

// collapsePaths схлопывает пересечения:
// если есть "2.7" — путь "2.7.9" не нужен.
func collapsePaths(paths []string) []string {
	uniq := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		uniq = append(uniq, p)
	}
	if len(uniq) <= 1 {
		return uniq
	}

	nlevel := func(p string) int { return 1 + strings.Count(p, ".") }

	sort.Slice(uniq, func(i, j int) bool {
		li, lj := nlevel(uniq[i]), nlevel(uniq[j])
		if li != lj {
			return li < lj
		}
		return uniq[i] < uniq[j]
	})

	out := make([]string, 0, len(uniq))
	for _, p := range uniq {
		redundant := false
		for _, keep := range out {
			if p == keep || strings.HasPrefix(p, keep+".") {
				redundant = true
				break
			}
		}
		if !redundant {
			out = append(out, p)
		}
	}
	return out
}

func (s *InfoBlockCollectionService) enrichInfoBlockResponses(infoBlocks []*models.InfoBlockResponse) {
	if len(infoBlocks) == 0 {
		return
	}

	templateIDs := make([]uint, 0)
	userIDs := make([]uint, 0)

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)

	resources := make([]contract.Resource, 0, len(infoBlocks))
	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}
		resources = append(resources, ib)

		if ib.TemplateID != nil {
			id := *ib.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}
		if ib.UserID != nil {
			id := *ib.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)

	var users map[uint]contract.User
	var templates map[uint]contract.Template
	var galleries map[uuid.UUID][]contract.Gallery

	go func() {
		defer wg.Done()
		if len(templateIDs) == 0 {
			return
		}
		var err error
		templates, err = s.api.Template.GetMapByIDs(templateIDs)
		if err != nil {
			logger.Errorf("[info_block][InfoBlockCollectionService][enrichInfoBlockResponses] Error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) == 0 {
			return
		}
		var err error
		users, err = s.api.User.GetMapByIDs(userIDs)
		if err != nil {
			logger.Errorf("[info_block][InfoBlockCollectionService][enrichInfoBlockResponses] Error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		galleries = s.api.Gallery.GetIndexesForResources(resources)
	}()

	wg.Wait()

	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}
		if galleries != nil {
			if g, ok := galleries[ib.UUID]; ok {
				ib.Galleries = g
			}
		}
		if templates != nil && ib.TemplateID != nil {
			ib.Template = templates[*ib.TemplateID]
		}
		if users != nil && ib.UserID != nil {
			ib.User = users[*ib.UserID]
		}
	}
}
