package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/axlle-com/blog/app/testutil"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupMenuRepo(t *testing.T) (repository.MenuItemRepository, *gorm.DB) {
	t.Helper()

	testDB, container := testutil.InitTestContainer(t)

	// isolate tests
	require.NoError(t, testDB.Exec(`TRUNCATE TABLE menu_items RESTART IDENTITY CASCADE;`).Error)

	return container.MenuItemRepo, testDB
}

func ptrUint(v uint) *uint { return &v }

func mustCreateMenuItem(t *testing.T, repo repository.MenuItemRepository, mi *models.MenuItem) *models.MenuItem {
	t.Helper()
	require.NoError(t, repo.Create(mi))
	require.NotZero(t, mi.ID)
	require.NotEmpty(t, mi.PathLtree)
	return mi
}

func mustGetMenuItem(t *testing.T, repo repository.MenuItemRepository, id uint) *models.MenuItem {
	t.Helper()
	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	return got
}

func buildMenuTree(t *testing.T, repo repository.MenuItemRepository, menuID uint) (root, child, grand *models.MenuItem) {
	t.Helper()

	root = mustCreateMenuItem(t, repo, &models.MenuItem{
		MenuID: menuID,
		Title:  "root",
		URL:    "/root",
		Sort:   1,
	})

	child = mustCreateMenuItem(t, repo, &models.MenuItem{
		MenuID:     menuID,
		MenuItemID: ptrUint(root.ID),
		Title:      "child",
		URL:        "/child",
		Sort:       2,
	})

	grand = mustCreateMenuItem(t, repo, &models.MenuItem{
		MenuID:     menuID,
		MenuItemID: ptrUint(child.ID),
		Title:      "grand",
		URL:        "/grand",
		Sort:       3,
	})

	return
}

func TestMenuItem_Create_RootChildGrandchild_PathLtree(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	root, child, grand := buildMenuTree(t, repo, 1)

	require.Equal(t, fmt.Sprintf("%d", root.ID), root.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root.ID, child.ID), child.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, child.ID, grand.ID), grand.PathLtree)
}

func TestMenuItem_GetDescendants_ReturnsAllSubtreeExceptSelf(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	root, child, grand := buildMenuTree(t, repo, 1)

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

func TestMenuItem_GetByFilter_ForNotMenuItemID_ExcludesNodeAndSubtree(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	root, child, grand := buildMenuTree(t, repo, 1)
	_ = root

	menuID := uint(1)
	forNot := child.ID

	items, err := repo.GetByFilter(nil, &models.MenuItemFilter{
		MenuID:           &menuID,
		ForNotMenuItemID: &forNot,
	})
	require.NoError(t, err)

	ids := map[uint]bool{}
	for _, it := range items {
		ids[it.ID] = true
	}

	// child и его поддерево не должны быть в выдаче
	require.False(t, ids[child.ID])
	require.False(t, ids[grand.ID])
}

func TestMenuItem_Update_SameParent_DoesNotChangePathLtree(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	root, child, _ := buildMenuTree(t, repo, 1)
	_ = root

	before := mustGetMenuItem(t, repo, child.ID)
	child.Title = "child-updated"

	require.NoError(t, repo.Update(child, nil))

	after := mustGetMenuItem(t, repo, child.ID)
	require.Equal(t, before.PathLtree, after.PathLtree)
	require.Equal(t, "child-updated", after.Title)
}

func TestMenuItem_Update_MoveToRoot_UpdatesSubtree_AndDoesNotTouchOtherMenu(t *testing.T) {
	repo, gdb := setupMenuRepo(t)

	// menu 1 tree
	root1, child1, grand1 := buildMenuTree(t, repo, 1)

	// menu 2 (must not be affected)
	root2, child2, _ := buildMenuTree(t, repo, 2)

	// move child1 to root (MenuItemID = nil)
	child1.MenuItemID = nil
	require.NoError(t, repo.Update(child1, nil))

	// verify menu1 paths
	child1DB := mustGetMenuItem(t, repo, child1.ID)
	grand1DB := mustGetMenuItem(t, repo, grand1.ID)

	require.Equal(t, fmt.Sprintf("%d", child1.ID), child1DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", child1.ID, grand1.ID), grand1DB.PathLtree)

	// verify menu2 untouched
	root2DB := mustGetMenuItem(t, repo, root2.ID)
	child2DB := mustGetMenuItem(t, repo, child2.ID)

	require.Equal(t, fmt.Sprintf("%d", root2.ID), root2DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root2.ID, child2.ID), child2DB.PathLtree)

	// sanity: menu1 root1 still exists and still has its original path
	root1DB := mustGetMenuItem(t, repo, root1.ID)
	require.Equal(t, fmt.Sprintf("%d", root1.ID), root1DB.PathLtree)

	// and nothing got "reset" across menus
	var cnt1, cnt2 int64
	require.NoError(t, gdb.Model(&models.MenuItem{}).Where("menu_id = 1").Count(&cnt1).Error)
	require.NoError(t, gdb.Model(&models.MenuItem{}).Where("menu_id = 2").Count(&cnt2).Error)
	require.GreaterOrEqual(t, cnt1, int64(3))
	require.GreaterOrEqual(t, cnt2, int64(3))
}

func TestMenuItem_Update_MoveUnderAnotherParent_UpdatesSubtree(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	// root
	root := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, Title: "root", URL: "/r"})

	// two children
	a := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, MenuItemID: ptrUint(root.ID), Title: "a", URL: "/a"})
	b := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, MenuItemID: ptrUint(root.ID), Title: "b", URL: "/b"})

	// grand under a
	ga := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, MenuItemID: ptrUint(a.ID), Title: "ga", URL: "/ga"})

	// move a under b
	a.MenuItemID = ptrUint(b.ID)
	require.NoError(t, repo.Update(a, nil))

	aDB := mustGetMenuItem(t, repo, a.ID)
	gaDB := mustGetMenuItem(t, repo, ga.ID)

	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, b.ID, a.ID), aDB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d.%d", root.ID, b.ID, a.ID, ga.ID), gaDB.PathLtree)
}

func TestMenuItem_Delete_Subtree(t *testing.T) {
	repo, gdb := setupMenuRepo(t)

	root, _, _ := buildMenuTree(t, repo, 1)

	require.NoError(t, repo.Delete(root))

	var cnt int64
	require.NoError(t, gdb.Model(&models.MenuItem{}).Where("menu_id = 1").Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestMenuItem_Delete_WithOnlyID_DeletesSubtree(t *testing.T) {
	repo, gdb := setupMenuRepo(t)

	root, _, _ := buildMenuTree(t, repo, 1)

	// delete only by ID (repo.Delete should fetch path_ltree itself)
	require.NoError(t, repo.Delete(&models.MenuItem{ID: root.ID}))

	var cnt int64
	require.NoError(t, gdb.Model(&models.MenuItem{}).Where("menu_id = 1").Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestMenuItem_GetByFilter_MenuIDAndIDs_NoLeakBetweenMenus(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	a := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, Title: "a", URL: "/a"})
	b := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 2, Title: "b", URL: "/b"})

	menuID := uint(1)
	items, err := repo.GetByFilter(nil, &models.MenuItemFilter{
		MenuID: &menuID,
		IDs:    []uint{a.ID, b.ID},
	})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, a.ID, items[0].ID)
	require.Equal(t, uint(1), items[0].MenuID)
}

func TestMenuItem_Update_MoveBetweenMenus_ShouldFail(t *testing.T) {
	repo, _ := setupMenuRepo(t)

	_, c1, _ := buildMenuTree(t, repo, 1)
	r2, _, _ := buildMenuTree(t, repo, 2)

	before := mustGetMenuItem(t, repo, c1.ID)

	// try to move c1 under r2 (different menu)
	c1.MenuItemID = ptrUint(r2.ID)
	err := repo.Update(c1, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot move node between menus")

	after := mustGetMenuItem(t, repo, before.ID)
	require.Equal(t, before.PathLtree, after.PathLtree)
	require.Equal(t, before.MenuID, after.MenuID)
}

func TestMenuItem_Update_InvalidOldPath_ShouldRefuse(t *testing.T) {
	repo, gdb := setupMenuRepo(t)

	root := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, Title: "r", URL: "/r"})
	child := mustCreateMenuItem(t, repo, &models.MenuItem{MenuID: 1, MenuItemID: ptrUint(root.ID), Title: "c", URL: "/c"})

	// Убедимся, что child реально != 0
	require.NotZero(t, child.ID)

	// делаем old path НЕвалидным: последний label НЕ равен child.ID
	// важно: это должен быть валидный ltree-литерал
	require.NoError(t, gdb.Exec(`UPDATE menu_items SET path_ltree = '1.bad' WHERE id = ?`, child.ID).Error)

	child.Title = "c2"
	err := repo.Update(child, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "old path_ltree is invalid")
}

func TestMenuItem_Update_ConcurrentWrite_WaitsAndSucceeds(t *testing.T) {
	repo, gdb := setupMenuRepo(t)

	_, child, _ := buildMenuTree(t, repo, 1)

	// Start tx1 and lock the row
	tx1 := gdb.Begin()
	require.NoError(t, tx1.Error)
	defer tx1.Rollback()

	var lockedID uint
	require.NoError(t, tx1.Raw(`SELECT id FROM menu_items WHERE id = ? FOR UPDATE`, child.ID).Scan(&lockedID).Error)
	require.Equal(t, child.ID, lockedID)

	// In parallel: Update should block until tx1 commits, then succeed.
	errCh := make(chan error, 1)
	durCh := make(chan time.Duration, 1)

	go func() {
		start := time.Now()
		child2 := &models.MenuItem{
			ID:         child.ID,
			MenuID:     child.MenuID,
			MenuItemID: child.MenuItemID,
			Title:      "child-concurrent",
			URL:        child.URL,
			Sort:       child.Sort,
		}
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
	// should have waited at least roughly "hold" (with some slack for scheduling)
	require.GreaterOrEqual(t, dur, 120*time.Millisecond)

	after := mustGetMenuItem(t, repo, child.ID)
	require.Equal(t, "child-concurrent", after.Title)
}
