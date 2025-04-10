package repository_test

import (
	"fmt"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	"testing"

	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	config.Config().SetTestENV()
	mPost.NewMigrator().Migrate()
	return db.GetDBTest()
}

func newTestRepo() repository.CategoryRepository {
	return repository.NewCategoryRepo().WithTx(setupTestDB())
}

func TestPathNotLikeQuery(t *testing.T) {
	// Инициализируем тестовую базу и репозиторий.
	db := setupTestDB()
	// Очищаем таблицу для чистоты теста.
	db.Exec("DELETE FROM post_categories")
	repo := newTestRepo()

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
	err = db.Where("path NOT LIKE ?", likePattern).Find(&results).Error
	assert.NoError(t, err)

	// Ожидаем, что в выборке будет только root2,
	// т.к. root1 и его дочерняя категория имеют пути, начинающиеся с root1.Path.
	assert.Equal(t, 1, len(results), fmt.Sprintf("Ожидалась 1 запись, получено %d", len(results)))
	assert.Equal(t, root2.ID, results[0].ID)
}

func TestPathLikeQuery(t *testing.T) {
	db := setupTestDB()
	db.Exec("DELETE FROM post_categories")
	repo := newTestRepo()

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
	err = db.Where("path LIKE ?", likePattern).Find(&results).Error
	assert.NoError(t, err)

	assert.Equal(t, 3, len(results))
}
