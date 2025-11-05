package db

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/blog/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

type seeder struct {
	postRepo          repository.PostRepository
	postService       *service.PostService
	categoryRepo      repository.CategoryRepository
	userProvider      user.UserProvider
	templateProvider  template.TemplateProvider
	infoBlockProvider apppPovider.InfoBlockProvider
	config            contract.Config
}

type PostSeedData struct {
	UserID             *uint    `json:"user_id"`
	Template           string   `json:"template"`
	PostCategoryID     *uint    `json:"post_category_id"`
	MetaTitle          *string  `json:"meta_title"`
	MetaDescription    *string  `json:"meta_description"`
	IsPublished        bool     `json:"is_published"`
	IsFavourites       bool     `json:"is_favourites"`
	HasComments        bool     `json:"has_comments"`
	ShowImagePost      bool     `json:"show_image_post"`
	ShowImageCategory  bool     `json:"show_image_category"`
	InSitemap          bool     `json:"in_sitemap"`
	IsMain             bool     `json:"is_main"`
	Media              *string  `json:"media"`
	Title              string   `json:"title"`
	TitleShort         *string  `json:"title_short"`
	DescriptionPreview *string  `json:"description_preview"`
	Description        *string  `json:"description"`
	ShowDate           bool     `json:"show_date"`
	DatePub            *string  `json:"date_pub"`
	DateEnd            *string  `json:"date_end"`
	Image              *string  `json:"image"`
	Hits               uint     `json:"hits"`
	Sort               int      `json:"sort"`
	Stars              float32  `json:"stars"`
	InfoBlocks         []string `json:"info_blocks"`
}

func NewSeeder(
	post repository.PostRepository,
	postService *service.PostService,
	category repository.CategoryRepository,
	user user.UserProvider,
	template template.TemplateProvider,
	infoBlockProvider apppPovider.InfoBlockProvider,
	cfg contract.Config,
) contract.Seeder {
	return &seeder{
		postRepo:          post,
		postService:       postService,
		categoryRepo:      category,
		userProvider:      user,
		templateProvider:  template,
		infoBlockProvider: infoBlockProvider,
		config:            cfg,
	}
}

func (s *seeder) Seed() error {
	return s.seedFromJSON("posts")
}

func (s *seeder) seedFromJSON(moduleName string) error {
	layout := s.config.Layout()
	seedPath := s.config.SrcFolderBuilder("db", layout, "seed", fmt.Sprintf("%s.json", moduleName))

	// Проверяем существование файла
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		logger.Infof("[blog][seeder][seedFromJSON] seed file not found: %s, skipping", seedPath)
		return nil
	}

	// Читаем JSON файл
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return err
	}

	var postsData []PostSeedData
	if err := json.Unmarshal(data, &postsData); err != nil {
		return err
	}

	idsUser := s.userProvider.GetAllIds()
	if len(idsUser) == 0 {
		return nil
	}

	for _, postData := range postsData {
		found, err := s.postService.FindByParam("title", postData.Title)
		var existingPost *models.Post
		if found != nil {
			existingPost = found
		}

		// Ищем шаблон по name и resource_name
		resourceName := "posts"
		var templateID *uint
		if postData.Template != "" {
			tpl, err := s.templateProvider.GetByNameAndResource(postData.Template, resourceName)
			if err != nil {
				logger.Errorf("[blog][seeder][seedFromJSON] template not found: name=%s, resource=%s, error=%v", postData.Template, resourceName, err)
				// Пропускаем этот пост, если шаблон не найден
				continue
			}
			id := tpl.GetID()
			templateID = &id
		}

		// Если пост уже существует, проверяем инфоблоки
		if existingPost != nil {
			// Получаем уже привязанные инфоблоки
			attachedInfoBlocks := s.infoBlockProvider.GetForResourceUUID(existingPost.UUID.String())
			attachedTitles := make(map[string]bool)
			for _, ib := range attachedInfoBlocks {
				attachedTitles[ib.GetTitle()] = true
			}

			// Привязываем недостающие инфоблоки
			if len(postData.InfoBlocks) > 0 {
				for _, infoBlockTitle := range postData.InfoBlocks {
					if infoBlockTitle == "" {
						continue
					}

					// Пропускаем, если инфоблок уже привязан
					if attachedTitles[infoBlockTitle] {
						continue
					}

					// Ищем инфоблок по title
					infoBlock, err := s.infoBlockProvider.FindByTitle(infoBlockTitle)
					if err != nil || infoBlock == nil {
						logger.Infof("[blog][seeder][seedFromJSON] info block with title='%s' not found for post '%s', skipping", infoBlockTitle, postData.Title)
						continue
					}

					// Привязываем инфоблок к посту
					_, err = s.infoBlockProvider.Attach(infoBlock.GetID(), existingPost.UUID.String())
					if err != nil {
						logger.Errorf("[blog][seeder][seedFromJSON] error attaching info block '%s' to post '%s': %v", infoBlockTitle, postData.Title, err)
						continue
					}
					logger.Infof("[blog][seeder][seedFromJSON] attached info block '%s' (ID=%d) to existing post '%s'", infoBlockTitle, infoBlock.GetID(), postData.Title)
				}
			}
			continue
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

		post := models.Post{
			UUID:               uuid.New(),
			TemplateID:         templateID,
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
			CreatedAt:          db.TimePtr(time.Now()),
			UpdatedAt:          db.TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		var userF contract.User
		if postData.UserID != nil {
			userF, _ = s.userProvider.GetByID(*postData.UserID)
		}

		createdPost, err := s.postService.Create(&post, userF)
		if err != nil {
			logger.Errorf("[blog][seeder][seedFromJSON] error creating post: %v", err)
			continue
		}

		// Привязываем инфоблоки по title, если они указаны
		if len(postData.InfoBlocks) > 0 {
			for _, infoBlockTitle := range postData.InfoBlocks {
				if infoBlockTitle == "" {
					continue
				}

				// Ищем инфоблок по title
				infoBlock, err := s.infoBlockProvider.FindByTitle(infoBlockTitle)
				if err != nil || infoBlock == nil {
					logger.Infof("[blog][seeder][seedFromJSON] info block with title='%s' not found for post '%s', skipping", infoBlockTitle, postData.Title)
					continue
				}

				// Привязываем инфоблок к посту
				_, err = s.infoBlockProvider.Attach(infoBlock.GetID(), createdPost.UUID.String())
				if err != nil {
					logger.Errorf("[blog][seeder][seedFromJSON] error attaching info block '%s' to post '%s': %v", infoBlockTitle, postData.Title, err)
					continue
				}
				logger.Infof("[blog][seeder][seedFromJSON] attached info block '%s' (ID=%d) to post '%s'", infoBlockTitle, infoBlock.GetID(), postData.Title)
			}
		}
	}

	logger.Infof("[blog][seeder][seedFromJSON] seeded %d posts from JSON", len(postsData))
	return nil
}

func (s *seeder) SeedTest(n int) error {
	err := s.categories(n)
	if err != nil {
		return err
	}

	return s.posts(n)
}

func (s *seeder) posts(n int) error {
	ids := s.templateProvider.GetAllIds()
	idsCategory, _ := s.categoryRepo.GetAllIds()
	idsUser := s.userProvider.GetAllIds()
	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomCategoryID := idsCategory[rand.Intn(len(idsCategory))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		post := models.Post{
			UUID:               uuid.New(),
			TemplateID:         &randomID,
			PostCategoryID:     &randomCategoryID,
			MetaTitle:          db.StrPtr(faker.Sentence()),
			MetaDescription:    db.StrPtr(faker.Sentence()),
			IsPublished:        db.RandBool(),
			IsFavourites:       db.RandBool(),
			HasComments:        db.RandBool(),
			ShowImagePost:      db.RandBool(),
			ShowImageCategory:  db.RandBool(),
			InSitemap:          db.RandBool(),
			Media:              db.StrPtr(faker.Word()),
			Title:              "TitlePost #" + strconv.Itoa(i),
			TitleShort:         db.StrPtr("TitlePostShort #" + strconv.Itoa(i)),
			DescriptionPreview: db.StrPtr(faker.Paragraph()),
			Description:        db.StrPtr(faker.Paragraph()),
			ShowDate:           db.RandBool(),
			DatePub:            db.ParseDate("02.01.2006"),
			DateEnd:            db.ParseDate("02.01.2006"),
			Image:              db.StrPtr("/public/img/404.svg"),
			Hits:               uint(rand.Intn(1000)),
			Sort:               rand.Intn(100),
			Stars:              rand.Float32() * 5,
			CreatedAt:          db.TimePtr(time.Now()),
			UpdatedAt:          db.TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		userF, _ := s.userProvider.GetByID(randomUserID)
		_, err := s.postService.Create(&post, userF)
		if err != nil {
			return err
		}
	}

	logger.Infof("[blog][seeder][SeedTest][posts] seeded %d posts", n)
	return nil
}

func (s *seeder) categories(n int) error {
	ids := s.templateProvider.GetAllIds()
	idsUser := s.userProvider.GetAllIds()

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

		randomID := ids[rand.Intn(len(ids))]
		postCategory := models.PostCategory{
			UUID:               uuid.New(),
			TemplateID:         &randomID,
			PostCategoryID:     randomCategoryID,
			UserID:             &randomUserID,
			MetaTitle:          db.StrPtr(faker.Sentence()),
			MetaDescription:    db.StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        db.IntToBoolPtr(),
			IsFavourites:       db.IntToBoolPtr(),
			InSitemap:          db.RandBool(),
			Title:              "TitleCategory #" + strconv.Itoa(i),
			TitleShort:         db.StrPtr("TitleCategoryShort #" + strconv.Itoa(i)),
			DescriptionPreview: db.StrPtr(faker.Paragraph()),
			Description:        db.StrPtr(faker.Paragraph()),
			Image:              db.StrPtr("/public/img/404.svg"),
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
