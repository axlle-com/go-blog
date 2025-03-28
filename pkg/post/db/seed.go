package db

import (
	. "github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/post/models"
	. "github.com/axlle-com/blog/pkg/post/repository"
	. "github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type seeder struct {
	postRepo         PostRepository
	postService      *PostService
	categoryRepo     CategoryRepository
	userProvider     user.UserProvider
	templateProvider template.TemplateProvider
}

func NewSeeder(
	post PostRepository,
	postService *PostService,
	category CategoryRepository,
	user user.UserProvider,
	template template.TemplateProvider,
) contracts.Seeder {
	return &seeder{
		postRepo:         post,
		postService:      postService,
		categoryRepo:     category,
		userProvider:     user,
		templateProvider: template,
	}
}

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	s.categories(n)
	s.posts(n)
}

func (s *seeder) posts(n int) {
	ids := s.templateProvider.GetAllIds()
	idsCategory, _ := s.categoryRepo.GetAllIds()
	idsUser := s.userProvider.GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomCategoryID := idsCategory[rand.Intn(len(idsCategory))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		post := Post{
			UUID:               uuid.New(),
			TemplateID:         &randomID,
			PostCategoryID:     &randomCategoryID,
			MetaTitle:          StrPtr(faker.Sentence()),
			MetaDescription:    StrPtr(faker.Sentence()),
			IsPublished:        RandBool(),
			IsFavourites:       RandBool(),
			HasComments:        RandBool(),
			ShowImagePost:      RandBool(),
			ShowImageCategory:  RandBool(),
			InSitemap:          RandBool(),
			Media:              StrPtr(faker.Word()),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			DescriptionPreview: StrPtr(faker.Paragraph()),
			Description:        StrPtr(faker.Paragraph()),
			ShowDate:           RandBool(),
			DatePub:            ParseDate("02.01.2006"),
			DateEnd:            ParseDate("02.01.2006"),
			Image:              StrPtr("/public/img/404.svg"),
			Hits:               uint(rand.Intn(1000)),
			Sort:               rand.Intn(100),
			Stars:              rand.Float32() * 5,
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		userF, _ := s.userProvider.GetByID(randomUserID)
		_, err := s.postService.Save(&post, userF)
		if err != nil {
			logger.Errorf("Failed to create userProvider %d: %v", i, err.Error())
		}
	}
	logger.Info("Database seeded Post successfully!")
}

func (s *seeder) categories(n int) {
	ids := s.templateProvider.GetAllIds()
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		postCategory := PostCategory{
			UUID:               uuid.New(),
			TemplateID:         &randomID,
			PostCategoryID:     UintPtr(rand.Intn(100)),
			MetaTitle:          StrPtr(faker.Sentence()),
			MetaDescription:    StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        IntToBoolPtr(),
			IsFavourites:       IntToBoolPtr(),
			InSitemap:          IntToBoolPtr(),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			DescriptionPreview: StrPtr(faker.Paragraph()),
			Description:        StrPtr(faker.Paragraph()),
			Image:              StrPtr("/public/img/404.svg"),
			Sort:               UintPtr(rand.Intn(100)),
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := s.categoryRepo.Create(&postCategory)
		if err != nil {
			logger.Errorf("Failed to create userProvider %d: %v", i, err.Error())
		}
	}
	logger.Info("Database seeded Post successfully!")
}
