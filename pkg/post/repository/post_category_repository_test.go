package repository

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/app/config"
	"github.com/axlle-com/blog/pkg/app/db"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	"github.com/bxcodec/faker/v3"
	"testing"

	postModel "github.com/axlle-com/blog/pkg/post/models" // импорт модели PostCategory
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	cfg := config.Config()
	cfg.SetTestENV()
	dbt := db.GetDBTest()
	mPost.NewMigrator().Rollback()
	mPost.NewMigrator().Migrate()
	return dbt
}

func newTestRepo() CategoryRepository {
	return &categoryRepository{
		db: setupTestDB(),
	}
}

func TestCreateRootCategory(t *testing.T) {
	repo := newTestRepo()

	root := &postModel.PostCategory{
		Title: "Root Category",
		Alias: faker.Username(),
		// Если PostCategoryID == nil, то это корневая категория.
	}
	err := repo.Create(root)
	assert.NoError(t, err)
	// Для первой корневой категории ожидаем left=1, right=2, level=0.
	assert.Equal(t, 1, root.LeftSet)
	assert.Equal(t, 2, root.RightSet)
	assert.Equal(t, 0, root.Level)
	assert.Nil(t, root.PostCategoryID)
}

func TestCreateChildCategory(t *testing.T) {
	repo := newTestRepo()

	// Создаем корневую категорию
	root := &postModel.PostCategory{
		Title: "Root Category",
		Alias: faker.Username(),
	}
	assert.NoError(t, repo.Create(root))

	// Создаем дочернюю категорию
	child := &postModel.PostCategory{
		Title:          "Child Category",
		Alias:          faker.Username(),
		PostCategoryID: &root.ID,
	}
	err := repo.Create(child)
	assert.NoError(t, err)

	// Ожидаем, что дочерняя категория получит left = 2, right = 3, level = 1.
	assert.Equal(t, 2, child.LeftSet)
	assert.Equal(t, 3, child.RightSet)
	assert.Equal(t, 1, child.Level)

	// Проверяем, что границы родительской категории скорректированы (right должно измениться).
	updatedRoot, err := repo.GetByID(root.ID)
	assert.NoError(t, err)
	// Если корень изначально имел right=2, после вставки дочернего ожидаем right = 4.
	assert.Equal(t, 1, updatedRoot.LeftSet)
	assert.Equal(t, 4, updatedRoot.RightSet)
}

func TestGetDescendants(t *testing.T) {
	repo := newTestRepo()

	// Строим дерево:
	// Root
	//   Child1
	//     Grandchild1
	//   Child2
	root := &postModel.PostCategory{Title: "Root", Alias: faker.Username()}
	assert.NoError(t, repo.Create(root))

	child1 := &postModel.PostCategory{Title: "Child1", PostCategoryID: &root.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(child1))

	grandchild1 := &postModel.PostCategory{Title: "Grandchild1", PostCategoryID: &child1.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(grandchild1))

	child2 := &postModel.PostCategory{Title: "Child2", PostCategoryID: &root.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(child2))

	// Перезагружаем root, чтобы получить актуальные LeftSet и RightSet
	updatedRoot, err := repo.GetByID(root.ID)
	assert.NoError(t, err)

	descendants, err := repo.GetDescendants(updatedRoot)
	assert.NoError(t, err)
	// Ожидаем 3 потомка: Child1, Grandchild1, Child2.
	assert.Len(t, descendants, 3)
	assert.Equal(t, child1.ID, descendants[0].ID)
	assert.Equal(t, grandchild1.ID, descendants[1].ID)
	assert.Equal(t, child2.ID, descendants[2].ID)

	// Проверяем получение потомков по ID
	descByID, err := repo.GetDescendantsByID(root.ID)
	assert.NoError(t, err)
	assert.Len(t, descByID, 3)
}

func TestUpdateCategory(t *testing.T) {
	repo := newTestRepo()

	// Создаем две корневые категории
	root1 := &postModel.PostCategory{Title: "Root1", Alias: faker.Username()}
	assert.NoError(t, repo.Create(root1))
	root2 := &postModel.PostCategory{Title: "Root2", Alias: faker.Username()}
	assert.NoError(t, repo.Create(root2))

	// Создаем дочернюю категорию под root1
	child := &postModel.PostCategory{
		Title:          "Child",
		Alias:          faker.Username(),
		PostCategoryID: &root1.ID,
	}
	assert.NoError(t, repo.Create(child))

	// Перемещаем child под root2
	oldChild := *child // копируем состояние до обновления
	child.PostCategoryID = &root2.ID
	child.Title = "Child moved"
	err := repo.Update(child, &oldChild)
	assert.NoError(t, err)

	// Проверяем, что у дочернего узла изменился родитель и уровень.
	updatedChild, err := repo.GetByID(child.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedChild.PostCategoryID)
	assert.Equal(t, root2.ID, *updatedChild.PostCategoryID)
	// Ожидаем, что дочерний узел теперь имеет level == 1.
	assert.Equal(t, 1, updatedChild.Level)
}

func TestDeleteCategory(t *testing.T) {
	repo := newTestRepo()

	// Строим дерево:
	// Root
	//   Child1
	//     Grandchild1
	//   Child2
	root := &postModel.PostCategory{Title: "Root", Alias: faker.Username()}
	assert.NoError(t, repo.Create(root))
	child1 := &postModel.PostCategory{Title: "Child1", PostCategoryID: &root.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(child1))
	grandchild1 := &postModel.PostCategory{Title: "Grandchild1", PostCategoryID: &child1.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(grandchild1))
	child2 := &postModel.PostCategory{Title: "Child2", PostCategoryID: &root.ID, Alias: faker.Username()}
	assert.NoError(t, repo.Create(child2))

	// Удаляем child1 — это должно удалить и его потомка grandchild1.
	assert.NoError(t, repo.Delete(child1))

	_, err := repo.GetByID(child1.ID)
	assert.Error(t, err)
	_, err = repo.GetByID(grandchild1.ID)
	assert.Error(t, err)

	// Проверяем, что child2 по-прежнему существует.
	fetchedChild2, err := repo.GetByID(child2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedChild2)
}

func TestWithPaginate(t *testing.T) {
	repo := newTestRepo()

	// Создаем 10 корневых категорий.
	for i := 1; i <= 10; i++ {
		cat := &postModel.PostCategory{
			Title: fmt.Sprintf("Category %d", i),
			Alias: faker.Username(),
		}
		assert.NoError(t, repo.Create(cat))
	}

	// Получаем первую страницу (page 1, pageSize 5).
	page1, err := repo.WithPaginate(1, 5)
	assert.NoError(t, err)
	assert.Len(t, page1, 5)

	// Получаем вторую страницу.
	page2, err := repo.WithPaginate(2, 5)
	assert.NoError(t, err)
	assert.Len(t, page2, 5)

	// Проверяем, что записи на страницах различны.
	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}
