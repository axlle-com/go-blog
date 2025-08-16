package repository_test

import (
	"fmt"
	"testing"

	"github.com/axlle-com/blog/pkg/blog/db/migrate"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
)

func setupTestDB() contracts.DB {
	cnf := config.Config()
	cnf.SetTestENV()
	//mPost.NewMigrator().Migrate()

	newDB, err := db.SetupDB(cnf)
	if err != nil {
		panic("db not initialized")
	}
	migrate.NewMigrator(newDB.PostgreSQL()).Migrate()
	return newDB
}

func newTestRepo(db contracts.DB) repository.CategoryRepository {
	return repository.NewCategoryRepo(db.PostgreSQL())
}

func TestPathNotLikeQuery(t *testing.T) {
	// Инициализируем тестовую базу и репозиторий.
	testDB := setupTestDB()
	// Очищаем таблицу для чистоты теста.
	testDB.PostgreSQL().Exec("DELETE FROM post_categories")
	repo := newTestRepo(testDB)

	// Создаем корневую категорию (будет иметь путь вида "/<ID>/")
	root1 := &models.PostCategory{
		Title: faker.Username(),
	}
	err := repo.Create(root1)
	assert.NoError(t, err)

	// Получаем актуальное значение root1.Path из базы.
	root1, err = repo.GetByID(root1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, root1.Path)

	// Создаем дочернюю категорию для root1, путь будет "/<root1.ID>/<child1.ID>/"
	child1 := &models.PostCategory{
		Title:          faker.Username(),
		PostCategoryID: &root1.ID,
	}
	err = repo.Create(child1)
	assert.NoError(t, err)

	// Создаем вторую корневую категорию, путь будет "/<ID>/"
	root2 := &models.PostCategory{
		Title: faker.Username(),
	}
	err = repo.Create(root2)
	assert.NoError(t, err)

	// Получаем актуальное значение root2.Path
	root2, err = repo.GetByID(root2.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, root2.Path)

	// Формируем шаблон для запроса: все записи, у которых путь начинается с root1.Path.
	likePattern := fmt.Sprintf("%s%%", root1.Path)

	// Выполняем запрос с NOT LIKE.
	var results []*models.PostCategory
	err = testDB.PostgreSQL().Where("path NOT LIKE ?", likePattern).Find(&results).Error
	assert.NoError(t, err)

	// Ожидаем, что в выборке будет только root2,
	// т.к. root1 и его дочерняя категория имеют пути, начинающиеся с root1.Path.
	assert.Equal(t, 1, len(results), fmt.Sprintf("Ожидалась 1 запись, получено %d", len(results)))
	assert.Equal(t, root2.ID, results[0].ID)
}

func TestPathLikeQuery(t *testing.T) {
	testDB := setupTestDB()
	testDB.PostgreSQL().Exec("DELETE FROM post_categories")
	repo := newTestRepo(testDB)

	// Создаем корневую категорию.
	root := &models.PostCategory{
		Title: faker.Username(),
	}
	err := repo.Create(root)
	assert.NoError(t, err)

	child1 := &models.PostCategory{
		Title:          faker.Username(),
		PostCategoryID: &root.ID,
	}
	err = repo.Create(child1)
	assert.NoError(t, err)

	child2 := &models.PostCategory{
		Title:          faker.Username(),
		PostCategoryID: &root.ID,
	}
	err = repo.Create(child2)
	assert.NoError(t, err)

	likePattern := fmt.Sprintf("%s%%", root.Path)

	var results []*models.PostCategory
	err = testDB.PostgreSQL().Where("path LIKE ?", likePattern).Find(&results).Error
	assert.NoError(t, err)

	assert.Equal(t, 3, len(results))
}
