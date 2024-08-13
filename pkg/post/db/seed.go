package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/post/models"
	templateRepo "github.com/axlle-com/blog/pkg/template/repository"
	userRepo "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedPosts(n int) {
	ids, _ := templateRepo.NewRepo().GetAllIds()
	idsCategory, _ := NewCategoryRepo().GetAllIds()
	idsUser, _ := userRepo.NewRepo().GetAllIds()
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
			DatePub:            ParseDate("02.01.2006"),
			DateEnd:            ParseDate("02.01.2006"),
			Image:              StrPtr(faker.Word()),
			Hits:               uint(rand.Intn(1000)),
			Sort:               rand.Intn(100),
			Stars:              rand.Float32() * 5,
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := NewPostRepo().Create(&post)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}

func SeedPostCategory(n int) {
	ids, _ := templateRepo.NewRepo().GetAllIds()
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		postCategory := PostCategory{
			TemplateID:         &randomID,
			PostCategoryID:     UintPtr(100),
			MetaTitle:          StrPtr(faker.Sentence()),
			MetaDescription:    StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        IntToBoolPtr(),
			IsFavourites:       IntToBoolPtr(),
			MakeWatermark:      IntToBoolPtr(),
			InSitemap:          IntToBoolPtr(),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			DescriptionPreview: StrPtr(faker.Paragraph()),
			Description:        StrPtr(faker.Paragraph()),
			Image:              StrPtr(faker.Word()),
			Sort:               UintPtr(100),
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := NewCategoryRepo().Create(&postCategory)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}
