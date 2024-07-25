package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/repository"
	categoryRepo "github.com/axlle-com/blog/pkg/post_category/repository"
	templateRepo "github.com/axlle-com/blog/pkg/template/repository"
	userRepo "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedPosts(n int) {
	ids, _ := templateRepo.NewRepository().GetAllIds()
	idsCategory, _ := categoryRepo.NewRepository().GetAllIds()
	idsUser, _ := userRepo.NewRepository().GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomCategoryID := ids[rand.Intn(len(idsCategory))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		post := Post{
			TemplateID:         &randomID,
			PostCategoryID:     &randomCategoryID,
			UserID:             &randomUserID,
			MetaTitle:          StrPtr(faker.Sentence()),
			MetaDescription:    StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        RandBool(),
			IsFavourites:       RandBool(),
			HasComments:        RandBool(),
			ShowImagePost:      RandBool(),
			ShowImageCategory:  RandBool(),
			MakeWatermark:      RandBool(),
			InSitemap:          RandBool(),
			Media:              StrPtr(faker.Word()),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			DescriptionPreview: StrPtr(faker.Paragraph()),
			Description:        StrPtr(faker.Paragraph()),
			ShowDate:           RandBool(),
			DatePub:            ParseDate(faker.Date()),
			DateEnd:            ParseDate(faker.Date()),
			Image:              StrPtr(faker.Word()),
			Hits:               uint(rand.Intn(1000)),
			Sort:               rand.Intn(100),
			Stars:              rand.Float32() * 5,
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := NewRepository().Create(&post)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}
