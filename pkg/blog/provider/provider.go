package provider

import (
	"net/url"
	"strconv"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
)

func NewBlogProvider(
	postService *service.PostService,
	postCollectionService *service.PostCollectionService,
	categoriesService *service.CategoriesService,
	tagCollectionService *service.TagCollectionService,
) contract.BlogProvider {
	return &provider{
		postService:           postService,
		postCollectionService: postCollectionService,
		categoriesService:     categoriesService,
		tagCollectionService:  tagCollectionService,
	}
}

type provider struct {
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoriesService     *service.CategoriesService
	tagCollectionService  *service.TagCollectionService
}

func (p *provider) GetPosts() []contract.Post {
	var collection []contract.Post
	posts, err := p.postCollectionService.GetAll()
	if err == nil {
		for _, post := range posts {
			collection = append(collection, post)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetPublishers(paginator contract.Paginator, filter contract.PublisherFilter) (collection []contract.Publisher, total int, err error) {
	// Создаем фильтры для всех источников
	postFilter := models.NewPostFilter()
	categoryFilter := models.NewCategoryFilterFilter()
	tagFilter := models.NewTagFilter()

	if filter != nil {
		if len(filter.GetUUIDs()) > 0 {
			postFilter.UUIDs = filter.GetUUIDs()
			categoryFilter.UUIDs = filter.GetUUIDs()
			tagFilter.UUIDs = filter.GetUUIDs()
		}
		if filter.GetQuery() != "" {
			query := filter.GetQuery()
			postFilter.Query = &query
			categoryFilter.Query = &query
			tagFilter.Query = &query
		}
	}

	pageSize := paginator.GetPageSize()
	page := paginator.GetPage()
	needed := pageSize
	globalOffset := (page - 1) * pageSize

	// Сначала получаем посты
	postsPaginator := paginator.Clone()
	posts, err := p.postCollectionService.WithPaginate(postsPaginator, postFilter)
	if err != nil {
		return
	}
	postsTotal := postsPaginator.GetTotal()

	// Добавляем посты в коллекцию
	for _, post := range posts {
		collection = append(collection, post)
		needed--
	}

	// Если набрали достаточно постов, получаем total из остальных источников и возвращаем
	if needed <= 0 {
		// Получаем total из категорий и тегов для общего total
		tempQuery := url.Values{}
		tempQuery.Set("page", "1")
		tempQuery.Set("pageSize", "1")
		tempPaginator := app.FromQuery(tempQuery)

		_, err = p.categoriesService.WithPaginate(tempPaginator, categoryFilter)
		if err != nil {
			return
		}
		categoriesTotal := tempPaginator.GetTotal()

		tempPaginator = app.FromQuery(tempQuery)
		_, err = p.tagCollectionService.WithPaginate(tempPaginator, tagFilter)
		if err != nil {
			return
		}
		tagsTotal := tempPaginator.GetTotal()

		total = postsTotal + categoriesTotal + tagsTotal
		return
	}

	// Если постов не хватило, добираем из категорий
	// Вычисляем offset для категорий:
	// - Если globalOffset < postsTotal: мы взяли посты, offset для категорий = 0 (начинаем с начала категорий)
	// - Если globalOffset >= postsTotal: постов не взяли, offset для категорий = globalOffset - postsTotal
	var categoriesOffset int
	if globalOffset < postsTotal {
		// Мы взяли посты, начинаем категории с начала
		categoriesOffset = 0
	} else {
		// Посты закончились, offset для категорий = сколько элементов мы пропустили в постах
		categoriesOffset = globalOffset - postsTotal
	}

	categoriesPage := (categoriesOffset / pageSize) + 1
	categoriesOffsetInPage := categoriesOffset % pageSize

	// Создаем пагинатор для категорий
	categoriesQuery := url.Values{}
	categoriesQuery.Set("page", strconv.Itoa(categoriesPage))
	categoriesQuery.Set("pageSize", strconv.Itoa(pageSize))
	categoriesPaginator := app.FromQuery(categoriesQuery)

	categories, err := p.categoriesService.WithPaginate(categoriesPaginator, categoryFilter)
	if err != nil {
		return
	}
	categoriesTotal := categoriesPaginator.GetTotal()

	// Добавляем категории в коллекцию
	start := categoriesOffsetInPage

	// Если offset указывает на элемент вне текущей страницы, берем следующую страницу
	if start >= len(categories) && categoriesTotal > 0 {
		// Текущая страница не содержит нужных элементов, берем следующую страницу
		categoriesPage++
		categoriesQuery.Set("page", strconv.Itoa(categoriesPage))
		categoriesPaginator = app.FromQuery(categoriesQuery)
		categories, err = p.categoriesService.WithPaginate(categoriesPaginator, categoryFilter)
		if err != nil {
			return
		}
		start = 0 // На новой странице начинаем с начала
	}

	end := start + needed
	if end > len(categories) {
		end = len(categories)
	}
	if start < len(categories) && end > start {
		for _, category := range categories[start:end] {
			collection = append(collection, category)
			needed--
		}
	}

	// Если набрали достаточно, получаем total из тегов и возвращаем
	if needed <= 0 {
		tempQuery := url.Values{}
		tempQuery.Set("page", "1")
		tempQuery.Set("pageSize", "1")
		tempPaginator := app.FromQuery(tempQuery)

		_, err = p.tagCollectionService.WithPaginate(tempPaginator, tagFilter)
		if err != nil {
			return
		}
		tagsTotal := tempPaginator.GetTotal()

		total = postsTotal + categoriesTotal + tagsTotal
		return
	}

	// Если категорий не хватило, добираем из тегов
	// Вычисляем offset для тегов:
	// - Если globalOffset < postsTotal: offset = 0 (начинаем теги с начала)
	// - Если globalOffset >= postsTotal: offset = (globalOffset - postsTotal) - categoriesTotal + сколько категорий взяли
	var tagsOffset int
	if globalOffset < postsTotal {
		// Мы взяли посты, начинаем теги с начала
		tagsOffset = 0
	} else {
		// Посты закончились, вычисляем offset относительно категорий
		categoriesOffsetForTags := globalOffset - postsTotal
		if categoriesOffsetForTags < categoriesTotal {
			// Мы взяли категории, начинаем теги с начала
			tagsOffset = 0
		} else {
			// Категории тоже закончились, offset для тегов = сколько элементов мы пропустили в категориях
			tagsOffset = categoriesOffsetForTags - categoriesTotal
		}
	}

	tagsPage := (tagsOffset / pageSize) + 1
	tagsOffsetInPage := tagsOffset % pageSize

	// Создаем пагинатор для тегов
	tagsQuery := url.Values{}
	tagsQuery.Set("page", strconv.Itoa(tagsPage))
	tagsQuery.Set("pageSize", strconv.Itoa(pageSize))
	tagsPaginator := app.FromQuery(tagsQuery)

	tags, err := p.tagCollectionService.WithPaginate(tagsPaginator, tagFilter)
	if err != nil {
		return
	}
	tagsTotal := tagsPaginator.GetTotal()

	// Добавляем теги в коллекцию
	start = tagsOffsetInPage
	end = start + needed
	if end > len(tags) {
		end = len(tags)
	}
	if start < len(tags) {
		for _, tag := range tags[start:end] {
			collection = append(collection, tag)
		}
	}

	// Общий total
	total = postsTotal + categoriesTotal + tagsTotal

	return
}
