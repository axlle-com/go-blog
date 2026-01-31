package db

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

type seeder struct {
	config          contract.Config
	disk            contract.DiskService
	db              contract.DB
	seedService     contract.SeedService
	api             *api.Api
	postRepo        repository.PostRepository
	postService     *service.PostService
	categoryRepo    repository.CategoryRepository
	categoryService *service.CategoryService
}

type InfoBlockSeedItem struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Sort     int    `json:"sort"`
	Position string `json:"position"`
}

type PostSeedData struct {
	UserID             *uint               `json:"user_id"`
	Template           string              `json:"template"`
	PostCategoryID     *uint               `json:"post_category_id"`
	CategoryAlias      *string             `json:"category_alias"`
	MetaTitle          *string             `json:"meta_title"`
	MetaDescription    *string             `json:"meta_description"`
	IsPublished        bool                `json:"is_published"`
	IsFavourites       bool                `json:"is_favourites"`
	HasComments        bool                `json:"has_comments"`
	ShowImagePost      bool                `json:"show_image_post"`
	ShowImageCategory  bool                `json:"show_image_category"`
	InSitemap          *bool               `json:"in_sitemap"`
	IsMain             bool                `json:"is_main"`
	Media              *string             `json:"media"`
	Title              string              `json:"title"`
	Alias              *string             `json:"alias"`
	TitleShort         *string             `json:"title_short"`
	DescriptionPreview *string             `json:"description_preview"`
	Description        *string             `json:"description"`
	ShowDate           bool                `json:"show_date"`
	DatePub            *string             `json:"date_pub"`
	DateEnd            *string             `json:"date_end"`
	Image              *string             `json:"image"`
	Hits               uint                `json:"hits"`
	Sort               int                 `json:"sort"`
	Stars              float32             `json:"stars"`
	InfoBlocks         []InfoBlockSeedItem `json:"info_blocks"`
}

type CategorySeedData struct {
	UserID          *uint               `json:"user_id"`
	Template        string              `json:"template"`
	IsPublished     *bool               `json:"is_published"`
	InSitemap       *bool               `json:"in_sitemap"`
	Alias           string              `json:"alias"`
	Title           string              `json:"title"`
	MetaTitle       string              `json:"meta_title,omitempty"`
	MetaDescription string              `json:"meta_description,omitempty"`
	InfoBlocks      []InfoBlockSeedItem `json:"info_blocks"`
}

func NewSeeder(
	cfg contract.Config,
	disk contract.DiskService,
	db contract.DB,
	seedService contract.SeedService,
	api *api.Api,
	post repository.PostRepository,
	postService *service.PostService,
	category repository.CategoryRepository,
	categoryService *service.CategoryService,
) contract.Seeder {
	return &seeder{
		config:          cfg,
		disk:            disk,
		db:              db,
		seedService:     seedService,
		api:             api,
		postRepo:        post,
		postService:     postService,
		categoryRepo:    category,
		categoryService: categoryService,
	}
}

func (s *seeder) Seed() error {
	if len(s.api.User.GetAllIds()) == 0 {
		return nil
	}

	if err := s.seedCategoriesFromJSON((&models.PostCategory{}).GetTable()); err != nil {
		return err
	}

	return s.seedFromJSON((&models.Post{}).GetTable())
}

func (s *seeder) SeedTest(n int) error {
	err := s.categories(n)
	if err != nil {
		return err
	}

	return s.posts(n)
}

func (s *seeder) seedFromJSON(moduleName string) error {
	files, _ := s.seedService.GetFiles(s.config.Layout(), moduleName)

	for name, seedPath := range files {
		data, err := s.disk.ReadFile(seedPath)
		if err != nil {
			return err
		}

		ok, err := s.seedService.IsApplied(name)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		var postsData []PostSeedData
		if err := json.Unmarshal(data, &postsData); err != nil {
			return err
		}

		for _, postData := range postsData {
			var found *models.Post
			if postData.Alias == nil {
				found, _ = s.postService.FindByParam("is_main", true)
			} else {
				found, _ = s.postService.FindByParam("alias", postData.Alias)
			}
			if found != nil {
				logger.Infof("[blog][seeder][seedFromJSON] post with title='%v' already exists, skipping", postData.Alias)

				continue
			}

			templateName := postData.Template

			if postData.CategoryAlias != nil {
				cat, err := s.categoryService.FindByParam("alias", *postData.CategoryAlias)
				if err != nil {
					logger.Errorf("[blog][seeder][seedFromJSON] category not found: alias=%s, error=%v", *postData.CategoryAlias, err)

					continue
				}
				id := cat.ID
				postData.PostCategoryID = &id
			}

			// Парсим даты
			var datePub, dateEnd *time.Time
			if postData.DatePub != nil {
				if parsed := db.ParseDate(*postData.DatePub); parsed != nil {
					datePub = parsed
				}
			}
			if postData.DateEnd != nil {
				if parsed := db.ParseDate(*postData.DateEnd); parsed != nil {
					dateEnd = parsed
				}
			}

			now := time.Now()

			post := models.Post{
				UUID:               uuid.New(),
				TemplateName:       templateName,
				PostCategoryID:     postData.PostCategoryID,
				MetaTitle:          postData.MetaTitle,
				MetaDescription:    postData.MetaDescription,
				IsPublished:        postData.IsPublished,
				IsFavourites:       postData.IsFavourites,
				HasComments:        postData.HasComments,
				ShowImagePost:      postData.ShowImagePost,
				ShowImageCategory:  postData.ShowImageCategory,
				InSitemap:          postData.InSitemap,
				IsMain:             postData.IsMain,
				Media:              postData.Media,
				Title:              postData.Title,
				TitleShort:         postData.TitleShort,
				DescriptionPreview: postData.DescriptionPreview,
				Description:        postData.Description,
				ShowDate:           postData.ShowDate,
				DatePub:            datePub,
				DateEnd:            dateEnd,
				Image:              postData.Image,
				Hits:               postData.Hits,
				Sort:               postData.Sort,
				Stars:              postData.Stars,
				CreatedAt:          db.TimePtr(now),
				UpdatedAt:          db.TimePtr(now),
				DeletedAt:          nil,
			}

			var userF contract.User
			if postData.UserID != nil {
				userF, _ = s.api.User.GetByID(*postData.UserID)
			}

			if postData.Alias != nil {
				post.Alias = *postData.Alias
			}

			createdPost, err := s.postService.Create(&post, userF)
			if err != nil {
				logger.Errorf("[blog][seeder][seedFromJSON] error creating post: %v", err)

				continue
			}

			// Привязываем инфоблоки по title, если они указаны
			if len(postData.InfoBlocks) > 0 {
				for i := range postData.InfoBlocks {
					infoBlockItem := &postData.InfoBlocks[i]
					if infoBlockItem.Title == "" {
						continue
					}

					// Ищем инфоблок по title
					infoBlock, err := s.api.InfoBlock.FindByTitle(infoBlockItem.Title)
					if err != nil || infoBlock == nil {
						logger.Infof("[blog][seeder][seedFromJSON] info block with title='%s' not found for post '%s', skipping", infoBlockItem.Title, postData.Title)

						continue
					}

					// Устанавливаем ID инфоблока
					infoBlockItem.ID = infoBlock.GetID()

					// Привязываем инфоблок к посту с указанием сортировки и позиции
					_, err = s.api.InfoBlock.CreateRelationFormBatch([]any{infoBlockItem}, createdPost.UUID.String())
					if err != nil {
						logger.Errorf("[blog][seeder][seedFromJSON] error attaching info block '%s' to post '%s': %v", infoBlockItem.Title, postData.Title, err)

						continue
					}

					logger.Infof("[blog][seeder][seedFromJSON] attached info block '%s' (ID=%d, Sort=%d, Position=%s) to post '%s'", infoBlockItem.Title, infoBlock.GetID(), infoBlockItem.Sort, infoBlockItem.Position, postData.Title)
				}
			}
		}

		if err := s.seedService.MarkApplied(name); err != nil {
			return err
		}

		logger.Infof("[blog][seeder][seedFromJSON] seeded %d posts from JSON", len(postsData))
	}

	return nil
}

func (s *seeder) seedCategoriesFromJSON(moduleName string) error {
	files, err := s.seedService.GetFiles(s.config.Layout(), moduleName)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logger.Infof("[blog][seeder][seedCategoriesFromJSON] seed files not found for module: %s, skipping", moduleName)
		return nil
	}

	for name, seedPath := range files {
		data, err := s.disk.ReadFile(seedPath)
		if err != nil {
			return err
		}

		ok, err := s.seedService.IsApplied(name)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		var categoriesData []CategorySeedData
		if err := json.Unmarshal(data, &categoriesData); err != nil {
			return err
		}

		for _, categoryData := range categoriesData {
			// Ищем существующую категорию по title
			allCategories, err := s.categoryRepo.GetAll()
			var existingCategory *models.PostCategory
			if err == nil {
				for _, cat := range allCategories {
					if cat.Title == categoryData.Title {
						existingCategory = cat
						break
					}
				}
			}

			// Если категория уже существует, пропускаем её
			if existingCategory != nil {
				logger.Infof("[blog][seeder][seedCategoriesFromJSON] category with title='%s' already exists, skipping", categoryData.Title)
				continue
			}

			templateName := categoryData.Template

			isPublished := true
			if categoryData.IsPublished != nil {
				isPublished = *categoryData.IsPublished
			}

			var inSitemap *bool
			if categoryData.InSitemap != nil {
				inSitemap = categoryData.InSitemap
			}

			category := models.PostCategory{
				UUID:            uuid.New(),
				TemplateName:    templateName,
				IsPublished:     &isPublished,
				InSitemap:       inSitemap,
				Alias:           categoryData.Alias,
				Title:           categoryData.Title,
				MetaTitle:       db.StrPtr(categoryData.MetaTitle),
				MetaDescription: db.StrPtr(categoryData.MetaDescription),
				CreatedAt:       db.TimePtr(time.Now()),
				UpdatedAt:       db.TimePtr(time.Now()),
				DeletedAt:       nil,
			}

			var userF contract.User
			if categoryData.UserID != nil {
				userF, _ = s.api.User.GetByID(*categoryData.UserID)
			}

			createdCategory, err := s.categoryService.Create(&category, userF)
			if err != nil {
				logger.Errorf("[blog][seeder][seedCategoriesFromJSON] error creating category: %v", err)
				continue
			}

			if len(categoryData.InfoBlocks) > 0 {
				for i := range categoryData.InfoBlocks {
					infoBlockItem := &categoryData.InfoBlocks[i]
					if infoBlockItem.Title == "" {
						continue
					}

					// Ищем инфоблок по title
					infoBlock, err := s.api.InfoBlock.FindByTitle(infoBlockItem.Title)
					if err != nil || infoBlock == nil {
						logger.Infof("[blog][seeder][seedCategoriesFromJSON] info block with title='%s' not found for category '%s', skipping", infoBlockItem.Title, categoryData.Title)
						continue
					}

					// Устанавливаем ID инфоблока
					infoBlockItem.ID = infoBlock.GetID()

					// Привязываем инфоблок к категории с указанием сортировки и позиции
					_, err = s.api.InfoBlock.CreateRelationFormBatch([]any{infoBlockItem}, createdCategory.UUID.String())
					if err != nil {
						logger.Errorf("[blog][seeder][seedCategoriesFromJSON] error attaching info block '%s' to category '%s': %v", infoBlockItem.Title, categoryData.Title, err)
						continue
					}
					logger.Infof("[blog][seeder][seedCategoriesFromJSON] attached info block '%s' (ID=%d, Sort=%d, Position=%s) to category '%s'", infoBlockItem.Title, infoBlock.GetID(), infoBlockItem.Sort, infoBlockItem.Position, categoryData.Title)
				}
			}
		}

		if err := s.seedService.MarkApplied(name); err != nil {
			return err
		}

		logger.Infof("[blog][seeder][seedCategoriesFromJSON] seeded %d categories from JSON (%s)", len(categoriesData), name)
	}

	return nil
}

func (s *seeder) posts(n int) error {
	templates := s.api.Template.GetAll()
	idsCategory, _ := s.categoryRepo.GetAllIds()
	idsUser := s.api.User.GetAllIds()
	for i := 1; i <= n; i++ {
		var templateName string
		if len(templates) > 0 {
			templateName = templates[rand.Intn(len(templates))].GetName()
		}
		randomCategoryID := idsCategory[rand.Intn(len(idsCategory))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		post := models.Post{
			UUID:               uuid.New(),
			TemplateName:       templateName,
			PostCategoryID:     &randomCategoryID,
			MetaTitle:          db.StrPtr(faker.Sentence()),
			MetaDescription:    db.StrPtr(faker.Sentence()),
			IsPublished:        db.RandBool(),
			IsFavourites:       db.RandBool(),
			HasComments:        db.RandBool(),
			ShowImagePost:      db.RandBool(),
			ShowImageCategory:  db.RandBool(),
			InSitemap:          db.IntToBoolPtr(),
			Media:              db.StrPtr(faker.Word()),
			Title:              "TitlePost #" + strconv.Itoa(i),
			TitleShort:         db.StrPtr("TitlePostShort #" + strconv.Itoa(i)),
			DescriptionPreview: db.StrPtr(faker.Paragraph()),
			Description:        db.StrPtr(faker.Paragraph()),
			ShowDate:           db.RandBool(),
			DatePub:            db.ParseDate("02.01.2006"),
			DateEnd:            db.ParseDate("02.01.2006"),
			Image:              db.StrPtr("/static/img/404.svg"),
			Hits:               uint(rand.Intn(1000)),
			Sort:               rand.Intn(100),
			Stars:              rand.Float32() * 5,
			CreatedAt:          db.TimePtr(time.Now()),
			UpdatedAt:          db.TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		userF, _ := s.api.User.GetByID(randomUserID)
		_, err := s.postService.Create(&post, userF)
		if err != nil {
			return err
		}
	}

	logger.Infof("[blog][seeder][SeedTest][posts] seeded %d posts", n)
	return nil
}

func (s *seeder) categories(n int) error {
	templates := s.api.Template.GetAll()
	idsUser := s.api.User.GetAllIds()

	for i := 1; i <= n; i++ {
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		idsCategory, _ := s.categoryRepo.GetAllIds()
		var randomCategoryID *uint
		if len(idsCategory) > 0 {
			randomCategoryID = &idsCategory[rand.Intn(len(idsCategory))]
			if rand.Intn(2) == 1 {
				randomCategoryID = nil
			}
		}

		var templateName string
		if len(templates) > 0 {
			templateName = templates[rand.Intn(len(templates))].GetName()
		}
		postCategory := models.PostCategory{
			UUID:               uuid.New(),
			TemplateName:       templateName,
			PostCategoryID:     randomCategoryID,
			UserID:             &randomUserID,
			MetaTitle:          db.StrPtr(faker.Sentence()),
			MetaDescription:    db.StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        db.IntToBoolPtr(),
			IsFavourites:       db.IntToBoolPtr(),
			InSitemap:          db.IntToBoolPtr(),
			Title:              "TitleCategory #" + strconv.Itoa(i),
			TitleShort:         db.StrPtr("TitleCategoryShort #" + strconv.Itoa(i)),
			DescriptionPreview: db.StrPtr(faker.Paragraph()),
			Description:        db.StrPtr(faker.Paragraph()),
			Image:              db.StrPtr("/static/img/404.svg"),
			Sort:               db.IntToUintPtr(rand.Intn(100)),
			CreatedAt:          db.TimePtr(time.Now()),
			UpdatedAt:          db.TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := s.categoryRepo.Create(&postCategory)
		if err != nil {
			return err
		}
	}

	logger.Infof("[blog][seeder][SeedTest][categories] seeded %d categories", n)
	return nil
}
