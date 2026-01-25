package repository_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	models2 "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/testutil"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupInfoBlockRepo(t *testing.T) (repository.InfoBlockRepository, *gorm.DB) {
	t.Helper()

	testDB, container := testutil.InitTestContainer(t)

	// isolate tests (не хардкодим имя таблицы)
	table := (&models.InfoBlock{}).GetTable()
	require.NoError(t, testDB.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table)).Error)

	// ВАЖНО: если в контейнере репозиторий называется иначе — поправь тут 1 строку.
	return container.InfoBlockRepo, testDB
}

func ptrUint(v uint) *uint       { return &v }
func ptrString(v string) *string { return &v }

func mustCreateInfoBlock(t *testing.T, repo repository.InfoBlockRepository, b *models.InfoBlock) *models.InfoBlock {
	t.Helper()

	if b.UUID == uuid.Nil {
		b.UUID = uuid.New()
	}
	if b.Title == "" {
		b.Title = "title-" + uuid.NewString()
	}
	if b.TemplateName == "" {
		b.TemplateName = "spring.info_blocks.default"
	}
	if b.UserID == nil {
		b.UserID = ptrUint(1)
	}

	require.NoError(t, repo.Create(b))
	require.NotZero(t, b.ID)
	require.NotEmpty(t, b.PathLtree)
	return b
}

func mustGetInfoBlock(t *testing.T, repo repository.InfoBlockRepository, id uint) *models.InfoBlock {
	t.Helper()
	got, err := repo.FindByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	return got
}

func buildInfoBlockTree(t *testing.T, repo repository.InfoBlockRepository, templateID, userID uint) (root, child, grand *models.InfoBlock) {
	t.Helper()
	templateIndex := fmt.Sprintf("spring.info_blocks.t%d", templateID)

	root = mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: templateIndex,
		UserID:       ptrUint(userID),
		Title:        "root-" + uuid.NewString(),
	})

	child = mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: templateIndex,
		UserID:       ptrUint(userID),
		InfoBlockID:  ptrUint(root.ID),
		Title:        "child-" + uuid.NewString(),
	})

	grand = mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: templateIndex,
		UserID:       ptrUint(userID),
		InfoBlockID:  ptrUint(child.ID),
		Title:        "grand-" + uuid.NewString(),
	})

	return
}

func TestInfoBlock_Create_RootChildGrandchild_PathLtree(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	root, child, grand := buildInfoBlockTree(t, repo, 1, 10)

	require.Equal(t, fmt.Sprintf("%d", root.ID), root.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root.ID, child.ID), child.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, child.ID, grand.ID), grand.PathLtree)
}

func TestInfoBlock_GetDescendants_ReturnsAllSubtreeExceptSelf(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	root, child, grand := buildInfoBlockTree(t, repo, 1, 10)

	desc, err := repo.GetDescendants(root)
	require.NoError(t, err)

	ids := map[uint]bool{}
	for _, it := range desc {
		ids[it.ID] = true
	}

	require.False(t, ids[root.ID])
	require.True(t, ids[child.ID])
	require.True(t, ids[grand.ID])
}

func TestInfoBlock_WithPaginate_FilterByTemplate_NoLeakBetweenTemplates(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	a := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		Title:        "a-" + uuid.NewString(),
	})
	_ = mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t2",
		UserID:       ptrUint(10),
		Title:        "b-" + uuid.NewString(),
	})

	p := models2.FromQuery(map[string][]string{})
	filter := &models.InfoBlockFilter{
		TemplateName: ptrString("spring.info_blocks.t1"),
	}
	items, err := repo.WithPaginate(p, filter)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(items), 1)

	// убеждаемся, что все из template=1
	for _, it := range items {
		require.NotEmpty(t, it.TemplateName)
		require.Equal(t, "spring.info_blocks.t1", it.TemplateName)
	}
	// и конкретно "a" присутствует
	foundA := false
	for _, it := range items {
		if it.ID == a.ID {
			foundA = true
			break
		}
	}
	require.True(t, foundA)
}

func TestInfoBlock_Update_SameParent_DoesNotChangePathLtree(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	_, child, _ := buildInfoBlockTree(t, repo, 1, 10)

	before := mustGetInfoBlock(t, repo, child.ID)

	updated := mustGetInfoBlock(t, repo, child.ID)
	updated.Title = "child-updated"

	require.NoError(t, repo.Update(updated, nil))

	after := mustGetInfoBlock(t, repo, child.ID)
	require.Equal(t, before.PathLtree, after.PathLtree)
	require.Equal(t, "child-updated", after.Title)
}

func TestInfoBlock_Update_MoveToRoot_UpdatesSubtree_AndDoesNotTouchOtherTree(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	// tree A
	_, child1, grand1 := buildInfoBlockTree(t, repo, 1, 10)

	// tree B (не должен быть затронут)
	root2, child2, _ := buildInfoBlockTree(t, repo, 2, 10)

	// move child1 to root
	toMove := mustGetInfoBlock(t, repo, child1.ID)
	toMove.InfoBlockID = nil
	require.NoError(t, repo.Update(toMove, nil))

	// verify A updated
	child1DB := mustGetInfoBlock(t, repo, child1.ID)
	grand1DB := mustGetInfoBlock(t, repo, grand1.ID)

	require.Equal(t, fmt.Sprintf("%d", child1.ID), child1DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", child1.ID, grand1.ID), grand1DB.PathLtree)

	// verify B untouched
	root2DB := mustGetInfoBlock(t, repo, root2.ID)
	child2DB := mustGetInfoBlock(t, repo, child2.ID)

	require.Equal(t, fmt.Sprintf("%d", root2.ID), root2DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root2.ID, child2.ID), child2DB.PathLtree)

	// sanity count (через Model, чтобы не попасть на soft-delete нюансы)
	table := (&models.InfoBlock{}).GetTable()
	var cntA, cntB int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t1", 10).Count(&cntA).Error)
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t2", 10).Count(&cntB).Error)
	require.GreaterOrEqual(t, cntA, int64(3))
	require.GreaterOrEqual(t, cntB, int64(3))
}

func TestInfoBlock_Update_MoveUnderAnotherParent_UpdatesSubtree(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	// root
	root := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		Title:        "root-" + uuid.NewString(),
	})

	// two children
	a := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		InfoBlockID:  ptrUint(root.ID),
		Title:        "a-" + uuid.NewString(),
	})
	b := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		InfoBlockID:  ptrUint(root.ID),
		Title:        "b-" + uuid.NewString(),
	})

	// grand under a
	ga := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		InfoBlockID:  ptrUint(a.ID),
		Title:        "ga-" + uuid.NewString(),
	})

	// move a under b
	aMove := mustGetInfoBlock(t, repo, a.ID)
	aMove.InfoBlockID = ptrUint(b.ID)
	require.NoError(t, repo.Update(aMove, nil))

	aDB := mustGetInfoBlock(t, repo, a.ID)
	gaDB := mustGetInfoBlock(t, repo, ga.ID)

	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, b.ID, a.ID), aDB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d.%d", root.ID, b.ID, a.ID, ga.ID), gaDB.PathLtree)
}

func TestInfoBlock_Update_CannotMoveUnderDescendant_ShouldRefuse(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	root, child, grand := buildInfoBlockTree(t, repo, 1, 10)

	// move root under grand => must fail
	rootMove := mustGetInfoBlock(t, repo, root.ID)
	rootMove.InfoBlockID = ptrUint(grand.ID)

	err := repo.Update(rootMove, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot move node under its descendant")

	// also move child under grand (direct descendant) => must fail
	childMove := mustGetInfoBlock(t, repo, child.ID)
	childMove.InfoBlockID = ptrUint(grand.ID)

	err = repo.Update(childMove, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot move node under its descendant")
}

func TestInfoBlock_Update_InvalidOldPath_ShouldRefuse(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	root := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		Title:        "r-" + uuid.NewString(),
	})
	child := mustCreateInfoBlock(t, repo, &models.InfoBlock{
		TemplateName: "spring.info_blocks.t1",
		UserID:       ptrUint(10),
		InfoBlockID:  ptrUint(root.ID),
		Title:        "c-" + uuid.NewString(),
	})

	// делаем old path НЕвалидным: последний label НЕ равен child.ID (но ltree литерал валидный)
	table := (&models.InfoBlock{}).GetTable()
	require.NoError(t, gdb.Exec(fmt.Sprintf(`UPDATE %s SET path_ltree = '%d.bad' WHERE id = ?`, table, root.ID), child.ID).Error)

	childUpd := mustGetInfoBlock(t, repo, child.ID)
	childUpd.Title = "c2-" + uuid.NewString()

	err := repo.Update(childUpd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "old path_ltree is invalid")
}

func TestInfoBlock_Delete_Subtree(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	root, _, _ := buildInfoBlockTree(t, repo, 1, 10)

	require.NoError(t, repo.Delete(root))

	table := (&models.InfoBlock{}).GetTable()
	var cnt int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t1", 10).Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestInfoBlock_Delete_WithOnlyID_DeletesSubtree(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	root, _, _ := buildInfoBlockTree(t, repo, 1, 10)

	// delete only by ID (repo.Delete should fetch path_ltree itself)
	require.NoError(t, repo.Delete(&models.InfoBlock{ID: root.ID}))

	table := (&models.InfoBlock{}).GetTable()
	var cnt int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t1", 10).Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestInfoBlock_DeleteByIDs_DeletesSubtrees(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	// tree A
	rootA, _, _ := buildInfoBlockTree(t, repo, 1, 10)

	// tree B
	rootB, childB, _ := buildInfoBlockTree(t, repo, 2, 10)

	// delete rootA subtree and childB subtree
	require.NoError(t, repo.DeleteByIDs([]uint{rootA.ID, childB.ID}))

	table := (&models.InfoBlock{}).GetTable()

	// A must be empty
	var cntA int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t1", 10).Count(&cntA).Error)
	require.Equal(t, int64(0), cntA)

	// B must still have rootB (childB removed and its subtree)
	var cntB int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.info_blocks.t2", 10).Count(&cntB).Error)
	require.Equal(t, int64(1), cntB)

	rootBDB := mustGetInfoBlock(t, repo, rootB.ID)
	require.Equal(t, fmt.Sprintf("%d", rootB.ID), rootBDB.PathLtree)
}

func TestInfoBlock_GetDescendantsByRoots_ReturnsUnionWithoutRoots(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	r1, c1, g1 := buildInfoBlockTree(t, repo, 1, 10)
	r2, c2, _ := buildInfoBlockTree(t, repo, 2, 10)

	desc, err := repo.GetDescendantsByRoots([]uint{r1.ID, r2.ID})
	require.NoError(t, err)

	ids := map[uint]bool{}
	for _, it := range desc {
		ids[it.ID] = true
	}

	// roots are excluded
	require.False(t, ids[r1.ID])
	require.False(t, ids[r2.ID])

	// descendants are included
	require.True(t, ids[c1.ID])
	require.True(t, ids[g1.ID])
	require.True(t, ids[c2.ID])
}

func TestInfoBlock_Update_ConcurrentWrite_WaitsAndSucceeds(t *testing.T) {
	repo, gdb := setupInfoBlockRepo(t)

	_, child, _ := buildInfoBlockTree(t, repo, 1, 10)

	// Start tx1 and lock the row
	tx1 := gdb.Begin()
	require.NoError(t, tx1.Error)
	defer tx1.Rollback()

	table := (&models.InfoBlock{}).GetTable()

	var lockedID uint
	require.NoError(t,
		tx1.Raw(fmt.Sprintf(`SELECT id FROM %s WHERE id = ? FOR UPDATE`, table), child.ID).
			Scan(&lockedID).Error,
	)
	require.Equal(t, child.ID, lockedID)

	// In parallel: Update should block until tx1 commits, then succeed.
	errCh := make(chan error, 1)
	durCh := make(chan time.Duration, 1)

	go func() {
		start := time.Now()

		child2 := mustGetInfoBlock(t, repo, child.ID)
		child2.Title = "child-concurrent-" + uuid.NewString()

		err := repo.Update(child2, nil)
		durCh <- time.Since(start)
		errCh <- err
	}()

	hold := 200 * time.Millisecond
	time.Sleep(hold) // keep lock for a bit
	require.NoError(t, tx1.Commit().Error)

	err := <-errCh
	dur := <-durCh

	require.NoError(t, err)
	require.GreaterOrEqual(t, dur, 120*time.Millisecond)

	after := mustGetInfoBlock(t, repo, child.ID)
	require.Contains(t, after.Title, "child-concurrent-")
}

func TestInfoBlock_Update_ParentToSelf_ShouldRefuse(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	root, _, _ := buildInfoBlockTree(t, repo, 1, 10)

	upd := mustGetInfoBlock(t, repo, root.ID)
	upd.InfoBlockID = ptrUint(root.ID)

	err := repo.Update(upd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot set parent to self")
}

func TestInfoBlock_Update_MoveUnderAnotherTemplateParent_ShouldSucceed(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	// дерево в template=1
	_, c1, g1 := buildInfoBlockTree(t, repo, 1, 10)

	// новый родитель в template=2
	r2, _, _ := buildInfoBlockTree(t, repo, 2, 10)

	// переносим c1 под r2
	c1Move := mustGetInfoBlock(t, repo, c1.ID)
	c1Move.InfoBlockID = ptrUint(r2.ID)

	require.NoError(t, repo.Update(c1Move, nil))

	c1DB := mustGetInfoBlock(t, repo, c1.ID)
	g1DB := mustGetInfoBlock(t, repo, g1.ID)

	require.Equal(t, fmt.Sprintf("%d.%d", r2.ID, c1.ID), c1DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d", r2.ID, c1.ID, g1.ID), g1DB.PathLtree)
}

func TestInfoBlock_GetDescendantsByPaths_Union_ExcludesRoots(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	// tree 1
	r1, c1, g1 := buildInfoBlockTree(t, repo, 1, 10)

	// tree 2
	r2, c2, _ := buildInfoBlockTree(t, repo, 2, 10)

	// tree 3 (не должен попасть)
	r3, c3, _ := buildInfoBlockTree(t, repo, 3, 10)

	paths := []string{
		r1.PathLtree,
		"   " + r2.PathLtree + "   ", // проверяем trim
		r1.PathLtree,                 // проверяем дубликаты
		"",                           // пустые игнорируются
	}

	desc, err := repo.GetDescendantsByPaths(paths)
	require.NoError(t, err)

	ids := map[uint]bool{}
	for _, it := range desc {
		ids[it.ID] = true
	}

	// roots excluded
	require.False(t, ids[r1.ID])
	require.False(t, ids[r2.ID])

	// descendants included
	require.True(t, ids[c1.ID])
	require.True(t, ids[g1.ID])
	require.True(t, ids[c2.ID])

	// чужая ветка не попала
	require.False(t, ids[r3.ID])
	require.False(t, ids[c3.ID])
}

func TestInfoBlock_GetSubtreesByPaths_Union_IncludesRoots(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	// tree 1
	r1, c1, g1 := buildInfoBlockTree(t, repo, 1, 10)

	// tree 2
	r2, c2, _ := buildInfoBlockTree(t, repo, 2, 10)

	// tree 3 (не должен попасть)
	r3, c3, _ := buildInfoBlockTree(t, repo, 3, 10)

	paths := []string{
		strings.TrimSpace(r1.PathLtree),
		r2.PathLtree,
		r2.PathLtree, // duplicate
	}

	all, err := repo.GetSubtreesByPaths(paths)
	require.NoError(t, err)

	ids := map[uint]bool{}
	for _, it := range all {
		ids[it.ID] = true
	}

	// roots included
	require.True(t, ids[r1.ID])
	require.True(t, ids[r2.ID])

	// descendants included
	require.True(t, ids[c1.ID])
	require.True(t, ids[g1.ID])
	require.True(t, ids[c2.ID])

	// чужая ветка не попала
	require.False(t, ids[r3.ID])
	require.False(t, ids[c3.ID])
}

func TestInfoBlock_GetDescendantsByPaths_EmptyInput_ReturnsEmpty(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	desc, err := repo.GetDescendantsByPaths([]string{})
	require.NoError(t, err)
	require.Len(t, desc, 0)

	desc, err = repo.GetDescendantsByPaths([]string{"", "   "})
	require.NoError(t, err)
	require.Len(t, desc, 0)
}

func TestInfoBlock_GetSubtreesByPaths_EmptyInput_ReturnsEmpty(t *testing.T) {
	repo, _ := setupInfoBlockRepo(t)

	all, err := repo.GetSubtreesByPaths([]string{})
	require.NoError(t, err)
	require.Len(t, all, 0)

	all, err = repo.GetSubtreesByPaths([]string{"", "   "})
	require.NoError(t, err)
	require.Len(t, all, 0)
}
