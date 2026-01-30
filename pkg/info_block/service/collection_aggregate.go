package service

import (
	"sort"
	"strings"
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
)

type CollectionAggregateService struct {
	api           *api.Api
	infoBlockRepo repository.InfoBlockRepository
	resourceRepo  repository.InfoBlockHasResourceRepository
}

func NewCollectionAggregateService(
	api *api.Api,
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
) *CollectionAggregateService {
	return &CollectionAggregateService{
		api:           api,
		infoBlockRepo: infoBlockRepo,
		resourceRepo:  resourceRepo,
	}
}

func (s *CollectionAggregateService) Aggregates(infoBlocks []*models.InfoBlock) []*models.InfoBlock {
	var templateNames []string
	var userIDs []uint

	templateNamesMap := make(map[string]bool)
	userIDsMap := make(map[uint]bool)

	for _, infoBlock := range infoBlocks {
		if infoBlock.TemplateName != "" && !templateNamesMap[infoBlock.TemplateName] {
			templateNames = append(templateNames, infoBlock.TemplateName)
			templateNamesMap[infoBlock.TemplateName] = true
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
	var templates map[string]contract.Template

	service.SafeGo(&wg, func() {
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Errorf("[info_block][CollectionService][Aggregates] Error: %v", err)
			}
		}
	})

	service.SafeGo(&wg, func() {
		if len(templateNames) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByNames(templateNames)
			if err != nil {
				logger.Errorf("[info_block][CollectionService][Aggregates] Error: %v", err)
			}
		}
	})

	wg.Wait()

	for _, infoBlock := range infoBlocks {
		if templates != nil && infoBlock.TemplateName != "" {
			infoBlock.Template = templates[infoBlock.TemplateName]
		}

		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return infoBlocks
}

func (s *CollectionAggregateService) AggregatesResponses(infoBlocks []*models.InfoBlockResponse) []*models.InfoBlockResponse {
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
			logger.Errorf("[info_block][CollectionService][AggregatesResponses] Error: %v", err)
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

	paths = s.collapsePaths(paths)
	if len(paths) == 0 {
		s.enrichInfoBlockResponses(infoBlocks)
		return infoBlocks
	}

	// 3) одним запросом получаем все узлы поддеревьев (включая корни)
	allNodes, err := s.infoBlockRepo.GetSubtreesByPaths(paths)
	if err != nil {
		logger.Errorf("[info_block][CollectionService][AggregatesResponses] Error: %v", err)
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
			ID:           n.ID,
			UUID:         n.UUID,
			TemplateName: n.TemplateName,
			InfoBlockID:  n.InfoBlockID,
			UserID:       n.UserID,
			Media:        n.Media,
			Title:        n.Title,
			Description:  n.Description,
			Image:        n.Image,
			PathLtree:    n.PathLtree,
			Sort:         n.Sort,
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

func (s *CollectionAggregateService) enrichInfoBlockResponses(infoBlocks []*models.InfoBlockResponse) {
	if len(infoBlocks) == 0 {
		return
	}

	templateNames := make([]string, 0)
	userIDs := make([]uint, 0)

	templateNamesMap := make(map[string]bool)
	userIDsMap := make(map[uint]bool)

	resources := make([]contract.Resource, 0, len(infoBlocks))
	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}
		resources = append(resources, ib)

		if ib.TemplateName != "" && !templateNamesMap[ib.TemplateName] {
			templateNames = append(templateNames, ib.TemplateName)
			templateNamesMap[ib.TemplateName] = true
		}

		if ib.UserID != nil {
			id := *ib.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
	}

	var users map[uint]contract.User
	var templates map[string]contract.Template
	var galleries map[uuid.UUID][]contract.Gallery

	if len(userIDs) != 0 {
		var err error
		users, err = s.api.User.GetMapByIDs(userIDs)
		if err != nil {
			logger.Errorf("[info_block][CollectionService][enrichInfoBlockResponses] Error: %v", err)
		}
	}

	galleries = s.api.Gallery.GetIndexesForResources(resources)

	if len(templateNames) != 0 {
		var err error
		templates, err = s.api.Template.GetMapByNames(templateNames)
		if err != nil {
			logger.Errorf("[info_block][CollectionService][enrichInfoBlockResponses] Error: %v", err)
		}
	}

	for _, ib := range infoBlocks {
		if ib == nil {
			continue
		}

		if galleries != nil {
			if g, ok := galleries[ib.UUID]; ok {
				ib.Galleries = g
			}
		}

		if templates != nil && ib.TemplateName != "" {
			ib.Template = templates[ib.TemplateName]
		}

		if users != nil && ib.UserID != nil {
			ib.User = users[*ib.UserID]
		}
	}
}

// collapsePaths схлопывает пересечения:
// если есть "2.7" — путь "2.7.9" не нужен.
func (s *CollectionAggregateService) collapsePaths(paths []string) []string {
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
