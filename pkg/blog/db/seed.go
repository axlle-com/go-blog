package db

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/repository"
	"github.com/axlle-com/blog/pkg/blog/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

type seeder struct {
	postRepo         repository.PostRepository
	postService      *service.PostService
	categoryRepo     repository.CategoryRepository
	userProvider     user.UserProvider
	templateProvider template.TemplateProvider
}

func NewSeeder(
	post repository.PostRepository,
	postService *service.PostService,
	category repository.CategoryRepository,
	user user.UserProvider,
	template template.TemplateProvider,
) contract.Seeder {
	return &seeder{
		postRepo:         post,
		postService:      postService,
		categoryRepo:     category,
		userProvider:     user,
		templateProvider: template,
	}
}

func (s *seeder) Seed() error {
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
	logger.Info("Database seeded Post successfully!")
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
	logger.Info("Database seeded Post successfully!")
	return nil
}
