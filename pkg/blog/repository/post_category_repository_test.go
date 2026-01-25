package repository_test

import (
	"fmt"
	"testing"
	"time"

	models2 "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/testutil"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func setupCategoryRepo(t *testing.T) (repository.CategoryRepository, *gorm.DB) {
	t.Helper()

	testDB, container := testutil.InitTestContainer(t)

	// isolate tests (не хардкодим имя таблицы)
	table := (&models.PostCategory{}).GetTable()
	require.NoError(t, testDB.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table)).Error)

	// ВАЖНО: если в контейнере репозиторий называется иначе — поправь тут 1 строку.
	return container.CategoryRepo, testDB
}

func ptrUint(v uint) *uint       { return &v }
func ptrBool(v bool) *bool       { return &v }
func ptrString(v string) *string { return &v }

func mustCreateCategory(t *testing.T, repo repository.CategoryRepository, c *models.PostCategory) *models.PostCategory {
	t.Helper()

	// минимально обязательные поля, чтобы не словить not-null/unique сюрпризы
	if c.UUID == uuid.Nil {
		c.UUID = uuid.New()
	}
	if c.Alias == "" {
		c.Alias = "alias-" + uuid.NewString()
	}
	if c.URL == "" {
		c.URL = "/url-" + uuid.NewString()
	}
	if c.Title == "" {
		c.Title = "title-" + uuid.NewString()
	}
	if c.TemplateName == "" {
		c.TemplateName = "spring.post_categories.default"
	}
	if c.UserID == nil {
		c.UserID = ptrUint(1)
	}
	if c.Sort == nil {
		c.Sort = ptrUint(0)
	}

	// jsonb поля: пусть будут валидные значения
	if len(c.GalleriesSnapshot) == 0 {
		c.GalleriesSnapshot = datatypes.JSON([]byte(`[]`))
	}
	if len(c.InfoBlocksSnapshot) == 0 {
		c.InfoBlocksSnapshot = datatypes.JSON([]byte(`[]`))
	}

	// булевы поля — пусть будут не nil (в Update они тоже сохраняются)
	if c.IsPublished == nil {
		c.IsPublished = ptrBool(true)
	}
	if c.InSitemap == nil {
		c.InSitemap = ptrBool(true)
	}
	if c.ShowImage == nil {
		c.ShowImage = ptrBool(true)
	}
	if c.IsFavourites == nil {
		c.IsFavourites = ptrBool(false)
	}

	require.NoError(t, repo.Create(c))
	require.NotZero(t, c.ID)
	require.NotEmpty(t, c.PathLtree)
	return c
}

func mustGetCategory(t *testing.T, repo repository.CategoryRepository, id uint) *models.PostCategory {
	t.Helper()
	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	return got
}

func buildCategoryTree(t *testing.T, repo repository.CategoryRepository, templateID, userID uint) (root, child, grand *models.PostCategory) {
	t.Helper()
	templateIndex := fmt.Sprintf("spring.post_categories.t%d", templateID)

	root = mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName: templateIndex,
		UserID:       ptrUint(userID),
		Title:        "root-" + uuid.NewString(),
		URL:          "/root-" + uuid.NewString(),
		Alias:        "root-" + uuid.NewString(),
		Sort:         ptrUint(1),
	})

	child = mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   templateIndex,
		UserID:         ptrUint(userID),
		PostCategoryID: ptrUint(root.ID),
		Title:          "child-" + uuid.NewString(),
		URL:            "/child-" + uuid.NewString(),
		Alias:          "child-" + uuid.NewString(),
		Sort:           ptrUint(2),
		MetaTitle:      ptrString("mt-child"),
	})

	grand = mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   templateIndex,
		UserID:         ptrUint(userID),
		PostCategoryID: ptrUint(child.ID),
		Title:          "grand-" + uuid.NewString(),
		URL:            "/grand-" + uuid.NewString(),
		Alias:          "grand-" + uuid.NewString(),
		Sort:           ptrUint(3),
	})

	return
}

func TestCategory_Create_RootChildGrandchild_PathLtree(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	root, child, grand := buildCategoryTree(t, repo, 1, 10)

	require.Equal(t, fmt.Sprintf("%d", root.ID), root.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root.ID, child.ID), child.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, child.ID, grand.ID), grand.PathLtree)
}

func TestCategory_GetDescendants_ReturnsAllSubtreeExceptSelf(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	root, child, grand := buildCategoryTree(t, repo, 1, 10)

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

func TestCategory_WithPaginate_FilterByTemplateAndUUIDs_NoLeakBetweenTemplates(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	a := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName: "spring.post_categories.t1",
		UserID:       ptrUint(10),
		Title:        "a-" + uuid.NewString(),
		URL:          "/a-" + uuid.NewString(),
		Alias:        "a-" + uuid.NewString(),
	})
	b := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName: "spring.post_categories.t2",
		UserID:       ptrUint(10),
		Title:        "b-" + uuid.NewString(),
		URL:          "/b-" + uuid.NewString(),
		Alias:        "b-" + uuid.NewString(),
	})

	p := models2.FromQuery(map[string][]string{})
	filter := &models.CategoryFilter{
		TemplateName: ptrString("spring.post_categories.t1"),
		UUIDs:        []uuid.UUID{a.UUID, b.UUID}, // b из другого template — не должен пролезть
	}
	items, err := repo.WithPaginate(p, filter)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, a.ID, items[0].ID)
	require.Equal(t, "spring.post_categories.t1", items[0].TemplateName)
}

func TestCategory_Update_SameParent_DoesNotChangePathLtree(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	_, child, _ := buildCategoryTree(t, repo, 1, 10)

	before := mustGetCategory(t, repo, child.ID)

	updated := mustGetCategory(t, repo, child.ID)
	updated.Title = "child-updated"
	// parent не меняем

	require.NoError(t, repo.Update(updated, nil))

	after := mustGetCategory(t, repo, child.ID)
	require.Equal(t, before.PathLtree, after.PathLtree)
	require.Equal(t, "child-updated", after.Title)
}

func TestCategory_Update_MoveToRoot_UpdatesSubtree_AndDoesNotTouchOtherScope(t *testing.T) {
	repo, gdb := setupCategoryRepo(t)

	// scope A: template=1 user=10
	root1, child1, grand1 := buildCategoryTree(t, repo, 1, 10)

	// scope B: template=2 user=10 (не должен быть затронут)
	root2, child2, _ := buildCategoryTree(t, repo, 2, 10)

	// move child1 to root
	toMove := mustGetCategory(t, repo, child1.ID)
	toMove.PostCategoryID = nil
	require.NoError(t, repo.Update(toMove, nil))

	// verify scope A updated
	child1DB := mustGetCategory(t, repo, child1.ID)
	grand1DB := mustGetCategory(t, repo, grand1.ID)

	require.Equal(t, fmt.Sprintf("%d", child1.ID), child1DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", child1.ID, grand1.ID), grand1DB.PathLtree)

	// verify scope B untouched
	root2DB := mustGetCategory(t, repo, root2.ID)
	child2DB := mustGetCategory(t, repo, child2.ID)

	require.Equal(t, fmt.Sprintf("%d", root2.ID), root2DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d", root2.ID, child2.ID), child2DB.PathLtree)

	// sanity: root1 still ok
	root1DB := mustGetCategory(t, repo, root1.ID)
	require.Equal(t, fmt.Sprintf("%d", root1.ID), root1DB.PathLtree)

	// and nothing got "reset" across scopes
	table := (&models.PostCategory{}).GetTable()
	var cntA, cntB int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.post_categories.t1", 10).Count(&cntA).Error)
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.post_categories.t2", 10).Count(&cntB).Error)
	require.GreaterOrEqual(t, cntA, int64(3))
	require.GreaterOrEqual(t, cntB, int64(3))
}

func TestCategory_Update_MoveUnderAnotherParent_UpdatesSubtree(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	// root
	root := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName: "spring.post_categories.t1",
		UserID:       ptrUint(10),
		Title:        "root-" + uuid.NewString(),
		URL:          "/r-" + uuid.NewString(),
		Alias:        "r-" + uuid.NewString(),
		Sort:         ptrUint(1),
	})

	// two children
	a := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   "spring.post_categories.t1",
		UserID:         ptrUint(10),
		PostCategoryID: ptrUint(root.ID),
		Title:          "a-" + uuid.NewString(),
		URL:            "/a-" + uuid.NewString(),
		Alias:          "a-" + uuid.NewString(),
		Sort:           ptrUint(2),
	})
	b := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   "spring.post_categories.t1",
		UserID:         ptrUint(10),
		PostCategoryID: ptrUint(root.ID),
		Title:          "b-" + uuid.NewString(),
		URL:            "/b-" + uuid.NewString(),
		Alias:          "b-" + uuid.NewString(),
		Sort:           ptrUint(3),
	})

	// grand under a
	ga := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   "spring.post_categories.t1",
		UserID:         ptrUint(10),
		PostCategoryID: ptrUint(a.ID),
		Title:          "ga-" + uuid.NewString(),
		URL:            "/ga-" + uuid.NewString(),
		Alias:          "ga-" + uuid.NewString(),
		Sort:           ptrUint(4),
	})

	// move a under b
	aMove := mustGetCategory(t, repo, a.ID)
	aMove.PostCategoryID = ptrUint(b.ID)
	require.NoError(t, repo.Update(aMove, nil))

	aDB := mustGetCategory(t, repo, a.ID)
	gaDB := mustGetCategory(t, repo, ga.ID)

	require.Equal(t, fmt.Sprintf("%d.%d.%d", root.ID, b.ID, a.ID), aDB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d.%d", root.ID, b.ID, a.ID, ga.ID), gaDB.PathLtree)
}

func TestCategory_Delete_Subtree(t *testing.T) {
	repo, gdb := setupCategoryRepo(t)

	root, _, _ := buildCategoryTree(t, repo, 1, 10)

	require.NoError(t, repo.Delete(root))

	table := (&models.PostCategory{}).GetTable()
	var cnt int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.post_categories.t1", 10).Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestCategory_Delete_WithOnlyID_DeletesSubtree(t *testing.T) {
	repo, gdb := setupCategoryRepo(t)

	root, _, _ := buildCategoryTree(t, repo, 1, 10)

	// delete only by ID (repo.Delete should fetch path_ltree itself)
	require.NoError(t, repo.Delete(&models.PostCategory{ID: root.ID}))

	table := (&models.PostCategory{}).GetTable()
	var cnt int64
	require.NoError(t, gdb.Table(table).Where("template_name = ? AND user_id = ?", "spring.post_categories.t1", 10).Count(&cnt).Error)
	require.Equal(t, int64(0), cnt)
}

func TestCategory_Update_MoveBetweenTemplates_ShouldSucceed(t *testing.T) {
	repo, _ := setupCategoryRepo(t)

	// дерево в template=1
	_, c1, g1 := buildCategoryTree(t, repo, 1, 10)

	// новый родитель в template=2 (раньше это считалось другим scope)
	r2, _, _ := buildCategoryTree(t, repo, 2, 10)

	// переносим c1 под r2
	c1Move := mustGetCategory(t, repo, c1.ID)
	c1Move.PostCategoryID = ptrUint(r2.ID)

	require.NoError(t, repo.Update(c1Move, nil))

	c1DB := mustGetCategory(t, repo, c1.ID)
	g1DB := mustGetCategory(t, repo, g1.ID)

	require.Equal(t, fmt.Sprintf("%d.%d", r2.ID, c1.ID), c1DB.PathLtree)
	require.Equal(t, fmt.Sprintf("%d.%d.%d", r2.ID, c1.ID, g1.ID), g1DB.PathLtree)
}

func TestCategory_Update_InvalidOldPath_ShouldRefuse(t *testing.T) {
	repo, gdb := setupCategoryRepo(t)

	root := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName: "spring.post_categories.t1",
		UserID:       ptrUint(10),
		Title:        "r-" + uuid.NewString(),
		URL:          "/r-" + uuid.NewString(),
		Alias:        "r-" + uuid.NewString(),
	})
	child := mustCreateCategory(t, repo, &models.PostCategory{
		TemplateName:   "spring.post_categories.t1",
		UserID:         ptrUint(10),
		PostCategoryID: ptrUint(root.ID),
		Title:          "c-" + uuid.NewString(),
		URL:            "/c-" + uuid.NewString(),
		Alias:          "c-" + uuid.NewString(),
	})

	require.NotZero(t, child.ID)

	// делаем old path НЕвалидным: последний label НЕ равен child.ID (но ltree литерал валидный)
	table := (&models.PostCategory{}).GetTable()
	require.NoError(t, gdb.Exec(fmt.Sprintf(`UPDATE %s SET path_ltree = '1.bad' WHERE id = ?`, table), child.ID).Error)

	childUpd := mustGetCategory(t, repo, child.ID)
	childUpd.Title = "c2-" + uuid.NewString()

	err := repo.Update(childUpd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "old path_ltree is invalid")
}

func TestCategory_Update_ConcurrentWrite_WaitsAndSucceeds(t *testing.T) {
	repo, gdb := setupCategoryRepo(t)

	_, child, _ := buildCategoryTree(t, repo, 1, 10)

	// Start tx1 and lock the row
	tx1 := gdb.Begin()
	require.NoError(t, tx1.Error)
	defer tx1.Rollback()

	table := (&models.PostCategory{}).GetTable()

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

		// Чтобы не затереть важные поля (Alias/URL/Title/и т.п.) берём актуальную запись и меняем Title
		child2 := mustGetCategory(t, repo, child.ID)
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
	// should have waited at least roughly "hold" (with some slack for scheduling)
	require.GreaterOrEqual(t, dur, 120*time.Millisecond)

	after := mustGetCategory(t, repo, child.ID)
	require.Contains(t, after.Title, "child-concurrent-")
}
