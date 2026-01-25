package repository

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InfoBlockRepository interface {
	WithTx(tx *gorm.DB) InfoBlockRepository
	Create(infoBlock *models.InfoBlock) error
	Update(new *models.InfoBlock, old *models.InfoBlock) error
	GetAll() ([]*models.InfoBlock, error)
	GetByIDs(ids []uint) ([]*models.InfoBlock, error)
	GetForResourceByFilter(filter *models.InfoBlockFilter) ([]*models.InfoBlockResponse, error)

	FindByID(id uint) (*models.InfoBlock, error)
	FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error)

	WithPaginate(p contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error)

	Delete(infoBlock *models.InfoBlock) error
	DeleteByIDs(ids []uint) (err error)

	GetRoots() ([]*models.InfoBlock, error)
	GetDescendants(infoBlock *models.InfoBlock) ([]*models.InfoBlock, error)
	GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error)

	// Потомки для нескольких корней
	GetDescendantsByRoots(rootIDs []uint) ([]*models.InfoBlock, error)
	// Потомки по массиву путей (исключая сами корни)
	GetDescendantsByPaths(paths []string) ([]*models.InfoBlock, error)

	// Поддеревья по массиву путей (включая сами корни)
	GetSubtreesByPaths(paths []string) ([]*models.InfoBlock, error)
}

type infoBlockRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewInfoBlockRepo(db *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: db}
}

func (r *infoBlockRepository) WithTx(tx *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: tx}
}

func (r *infoBlockRepository) Create(infoBlock *models.InfoBlock) error {
	if infoBlock == nil {
		return fmt.Errorf("infoBlock is nil")
	}

	infoBlock.Creating()

	// root
	if infoBlock.InfoBlockID == nil || *infoBlock.InfoBlockID == 0 {
		if err := r.db.Create(infoBlock).Error; err != nil {
			return err
		}

		infoBlock.PathLtree = fmt.Sprintf("%d", infoBlock.ID)

		return r.db.Model(infoBlock).UpdateColumn("path_ltree", infoBlock.PathLtree).Error
	}

	// parent
	var parent models.InfoBlock
	if err := r.db.First(&parent, *infoBlock.InfoBlockID).Error; err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}

	if parent.PathLtree == "" {
		return fmt.Errorf("parent path_ltree is empty (id=%d)", parent.ID)
	}

	if err := r.db.Create(infoBlock).Error; err != nil {
		return err
	}

	infoBlock.PathLtree = fmt.Sprintf("%s.%d", parent.PathLtree, infoBlock.ID)

	return r.db.Model(infoBlock).UpdateColumn("path_ltree", infoBlock.PathLtree).Error
}

func (r *infoBlockRepository) FindByID(id uint) (*models.InfoBlock, error) {
	var model models.InfoBlock
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *infoBlockRepository) FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error) {
	var model models.InfoBlock
	query := r.db.Model(&model)

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.Title != nil {
		query = query.Where("title = ?", *filter.Title)
	}

	if filter.TemplateName != nil && *filter.TemplateName != "" {
		query = query.Where("template_name = ?", *filter.TemplateName)
	}

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if err := query.First(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *infoBlockRepository) WithPaginate(p contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	var total int64

	q := r.db.Model(&models.InfoBlock{})

	if filter != nil {
		if filter.ID != nil {
			q = q.Where("id = ?", *filter.ID)
		}
		if filter.TemplateName != nil && *filter.TemplateName != "" {
			q = q.Where("template_name = ?", *filter.TemplateName)
		}
		if filter.UserID != nil {
			q = q.Where("user_id = ?", *filter.UserID)
		}
		if filter.Title != nil && *filter.Title != "" {
			q = q.Where("title ILIKE ?", "%"+*filter.Title+"%")
		}
		if len(filter.UUIDs) > 0 {
			q = q.Where("uuid IN ?", filter.UUIDs)
		}
		if len(filter.IDs) > 0 {
			q = q.Where("id IN ?", filter.IDs)
		}
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	err := q.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order("id ASC").
		Find(&infoBlocks).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))

	return infoBlocks, nil
}

func (r *infoBlockRepository) Update(new *models.InfoBlock, _ *models.InfoBlock) error {
	if new == nil {
		return fmt.Errorf("new is nil")
	}

	new.Updating()

	return r.db.Transaction(func(db *gorm.DB) error {
		// 1) читаем old под блокировкой
		var old models.InfoBlock
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&old, new.ID).Error; err != nil {
			return err
		}

		if err := validatePathEndsWithID(old.PathLtree, old.ID); err != nil {
			return fmt.Errorf("refuse to update because old path_ltree is invalid: %w", err)
		}

		// 2) определяем oldParent/newParent
		oldParent, newParent := uint(0), uint(0)
		if old.InfoBlockID != nil {
			oldParent = *old.InfoBlockID
		}
		if new.InfoBlockID != nil {
			newParent = *new.InfoBlockID
		}

		// 2.0) защита: parent=self
		if newParent != 0 && newParent == old.ID {
			return fmt.Errorf("cannot set parent to self: id=%d", old.ID)
		}

		new.TemplateName = old.TemplateName
		new.UserID = old.UserID

		// 2.1) если родитель не менялся — обычный save, path_ltree оставляем прежним
		if oldParent == newParent {
			new.PathLtree = old.PathLtree
			return db.Select(new.Fields()).Save(new).Error
		}

		// 3) валидируем нового родителя (если перенос не в root)
		var newParentInfoBlock models.InfoBlock
		if newParent != 0 {
			if err := db.Select("id", "path_ltree").
				First(&newParentInfoBlock, newParent).Error; err != nil {
				return fmt.Errorf("new parent not found: %w", err)
			}

			if newParentInfoBlock.PathLtree == "" {
				return fmt.Errorf("new parent has empty path_ltree (id=%d)", newParentInfoBlock.ID)
			}

			// защита: нельзя завернуть в своего потомка
			oldPrefix := old.PathLtree
			descPrefix := oldPrefix + "."
			if newParentInfoBlock.PathLtree == oldPrefix || strings.HasPrefix(newParentInfoBlock.PathLtree, descPrefix) {
				return fmt.Errorf(
					"cannot move node under its descendant: node=%d new_parent=%d",
					old.ID, newParentInfoBlock.ID,
				)
			}
		}

		oldPathLtree := old.PathLtree

		var newPathLtree string
		if newParent == 0 {
			newPathLtree = fmt.Sprintf("%d", new.ID)
		} else {
			newPathLtree = fmt.Sprintf("%s.%d", newParentInfoBlock.PathLtree, new.ID)
		}

		// 4) обновляем потомков (кроме корня поддерева)
		err := db.Exec(fmt.Sprintf(`
			UPDATE %s
			SET path_ltree = ?::ltree || subpath(path_ltree, nlevel(?::ltree))
			WHERE path_ltree <@ ?::ltree
			  AND path_ltree <> ?::ltree
		`, new.GetTable()),
			newPathLtree, oldPathLtree,
			oldPathLtree, oldPathLtree,
		).Error
		if err != nil {
			return fmt.Errorf("failed to update subtree paths: %w", err)
		}

		// 5) обновляем сам узел
		new.PathLtree = newPathLtree
		return db.Select(new.Fields()).Save(new).Error
	})
}

func (r *infoBlockRepository) Delete(infoBlock *models.InfoBlock) error {
	if infoBlock == nil {
		return fmt.Errorf("infoBlock is nil")
	}

	if !infoBlock.Deleting() {
		return errors.New("deletion errors occurred")
	}

	// если PathLtree не заполнен — подстрахуемся и возьмём из БД
	path := infoBlock.PathLtree
	if path == "" && infoBlock.ID != 0 {
		var tmp models.InfoBlock
		if err := r.db.Select("path_ltree").First(&tmp, infoBlock.ID).Error; err != nil {
			return err
		}
		path = tmp.PathLtree
	}

	if path == "" {
		return fmt.Errorf("empty path_ltree for delete (id=%d)", infoBlock.ID)
	}

	// узел + потомки
	return r.db.Where("path_ltree <@ ?::ltree", path).
		Delete(&models.InfoBlock{}).Error
}

func (r *infoBlockRepository) GetAll() ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	if err := r.db.Order("id ASC").Find(&infoBlocks).Error; err != nil {
		return nil, err
	}

	return infoBlocks, nil
}

func (r *infoBlockRepository) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	if len(ids) == 0 {
		return []*models.InfoBlock{}, nil
	}

	if err := r.db.Where("id IN ?", ids).Find(&infoBlocks).Error; err != nil {
		return nil, err
	}
	return infoBlocks, nil
}

// DeleteByIDs удаляет не только сами узлы, но и их поддеревья (ltree).
// Иначе легко получить "сирот" (потомков без родителя).
func (r *infoBlockRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	// 1) грузим узлы (как минимум id/path_ltree)
	type row struct {
		ID        uint
		PathLtree string
	}
	var rows []row
	if err := r.db.Model(&models.InfoBlock{}).
		Select("id, path_ltree").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	// 2) если у кого-то пустой path_ltree — дотянем его отдельным запросом
	missing := make([]uint, 0)
	for _, rt := range rows {
		if rt.PathLtree == "" {
			missing = append(missing, rt.ID)
		}
	}
	if len(missing) > 0 {
		var fix []row
		if err := r.db.Model(&models.InfoBlock{}).
			Select("id, path_ltree").
			Where("id IN ?", missing).
			Find(&fix).Error; err != nil {
			return err
		}

		// мержим
		mp := make(map[uint]string, len(fix))
		for _, f := range fix {
			mp[f.ID] = f.PathLtree
		}
		for i := range rows {
			if rows[i].PathLtree == "" {
				rows[i].PathLtree = mp[rows[i].ID]
			}
		}
	}

	paths := make([]string, 0, len(rows))
	for _, rt := range rows {
		if rt.PathLtree != "" {
			paths = append(paths, rt.PathLtree)
		}
	}
	if len(paths) == 0 {
		return fmt.Errorf("empty path_ltree for delete by ids")
	}

	// 3) собираем OR по всем корням удаления: path_ltree <@ rootPath
	q := r.db.Model(&models.InfoBlock{}).Where("1=0")
	for _, p := range paths {
		q = q.Or("path_ltree <@ ?::ltree", p)
	}

	// 4) удаляем все поддеревья одним запросом
	return q.Delete(&models.InfoBlock{}).Error
}

func (r *infoBlockRepository) GetForResourceByFilter(filter *models.InfoBlockFilter) ([]*models.InfoBlockResponse, error) {
	var infoBlocks []*models.InfoBlockResponse
	query := r.db.
		Joins("inner join info_block_has_resources as r on info_blocks.id = r.info_block_id").
		Select("info_blocks.*", "r.id as relation_id", "r.sort as sort", "r.position as position", "r.resource_uuid as resource_uuid")

	if filter != nil {
		if filter.ResourceUUID != nil {
			query = query.Where("r.resource_uuid = ?", filter.ResourceUUID)
		}

		if filter.RelationID != nil {
			query = query.Where("r.id = ?", filter.RelationID)
		}

		if filter.RelationIDs != nil && len(filter.RelationIDs) > 0 {
			query = query.Where("r.id IN ?", filter.RelationIDs)
		}

		if filter.ID != nil {
			query = query.Where("info_blocks.id = ?", filter.ID)
		}

		if filter.IDs != nil && len(filter.IDs) > 0 {
			query = query.Where("info_blocks.id IN ?", filter.IDs)
		}

		if filter.UUIDs != nil && len(filter.UUIDs) > 0 {
			query = query.Where("info_blocks.uuid IN ?", filter.UUIDs)
		}
	}

	query = query.Order("r.sort ASC").
		Order("r.position ASC").
		Order("r.id DESC").
		Model(&models.InfoBlock{})

	err := query.Find(&infoBlocks).Error

	return infoBlocks, err
}

func (r *infoBlockRepository) GetRoots() ([]*models.InfoBlock, error) {
	var roots []*models.InfoBlock
	err := r.db.Where("info_block_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *infoBlockRepository) GetDescendants(infoBlock *models.InfoBlock) ([]*models.InfoBlock, error) {
	var descendants []*models.InfoBlock
	if infoBlock == nil || infoBlock.ID == 0 || infoBlock.PathLtree == "" {
		return []*models.InfoBlock{}, nil
	}

	err := r.db.
		Where("path_ltree <@ ?::ltree AND id <> ?", infoBlock.PathLtree, infoBlock.ID).
		Order("id ASC").
		Find(&descendants).Error

	return descendants, err
}

func (r *infoBlockRepository) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	// "все элементы, которые НЕ являются потомками parent (и не сам parent)"
	var items []*models.InfoBlock
	if parent == nil || parent.ID == 0 || parent.PathLtree == "" {
		return []*models.InfoBlock{}, nil
	}

	err := r.db.
		Where("NOT (path_ltree <@ ?::ltree) AND id <> ?", parent.PathLtree, parent.ID).
		Order("id ASC").
		Find(&items).Error

	return items, err
}

func (r *infoBlockRepository) GetDescendantsByRoots(rootIDs []uint) ([]*models.InfoBlock, error) {
	if len(rootIDs) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// берём пути корней
	type row struct {
		ID        uint
		PathLtree string
	}
	var roots []row
	if err := r.db.Model(&models.InfoBlock{}).
		Select("id, path_ltree").
		Where("id IN ?", rootIDs).
		Find(&roots).Error; err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(roots))
	for _, rt := range roots {
		if rt.PathLtree != "" {
			paths = append(paths, rt.PathLtree)
		}
	}
	if len(paths) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// ВАЖНО: <@ включает сам root, поэтому исключаем root сразу в каждой OR-ветке
	q := r.db.Model(&models.InfoBlock{}).Where("1=0")
	for _, p := range paths {
		q = q.Or("(path_ltree <@ ?::ltree AND path_ltree <> ?::ltree)", p, p)
	}

	var descendants []*models.InfoBlock
	err := q.Order("id ASC").Find(&descendants).Error
	if err != nil {
		return nil, err
	}

	return descendants, nil
}

func (r *infoBlockRepository) GetDescendantsByPaths(paths []string) ([]*models.InfoBlock, error) {
	if len(paths) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// чистим пустые и дубли
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

	if len(uniq) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// OR по всем путям: (path_ltree <@ p AND path_ltree <> p)
	q := r.db.Model(&models.InfoBlock{}).Where("1=0")
	for _, p := range uniq {
		q = q.Or("(path_ltree <@ ?::ltree AND path_ltree <> ?::ltree)", p, p)
	}

	var out []*models.InfoBlock
	if err := q.Order("id ASC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *infoBlockRepository) GetSubtreesByPaths(paths []string) ([]*models.InfoBlock, error) {
	if len(paths) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// чистим пустые и дубли
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
	if len(uniq) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// OR по всем путям: path_ltree <@ p  (включая корень)
	q := r.db.Model(&models.InfoBlock{}).Where("1=0")
	for _, p := range uniq {
		q = q.Or("path_ltree <@ ?::ltree", p)
	}

	var out []*models.InfoBlock
	if err := q.Order("id ASC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func validatePathEndsWithID(path string, id uint) error {
	if path == "" {
		return fmt.Errorf("empty path_ltree")
	}

	last := path
	if i := strings.LastIndex(path, "."); i >= 0 {
		last = path[i+1:]
	}

	want := fmt.Sprintf("%d", id)
	if last != want {
		return fmt.Errorf("bad path_ltree=%q: last label=%q, expected id=%q", path, last, want)
	}

	return nil
}

func collapsePaths(paths []string) []string {
	// 1) trim + unique
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

	// 2) сортируем по длине (кол-во сегментов), потом лексикографически
	nlevel := func(p string) int {
		// "2.7.9" => 3
		return 1 + strings.Count(p, ".")
	}
	sort.Slice(uniq, func(i, j int) bool {
		li, lj := nlevel(uniq[i]), nlevel(uniq[j])
		if li != lj {
			return li < lj
		}
		return uniq[i] < uniq[j]
	})

	// 3) оставляем только верхние пути
	out := make([]string, 0, len(uniq))
	for _, p := range uniq {
		redundant := false
		for _, keep := range out {
			// p внутри keep? (keep == "2.7", p == "2.7.9")
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
