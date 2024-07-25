package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post_category/repository"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedPostCategory(n int) {
	ids, _ := repository.NewRepository().GetAllIds()
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

		err := NewRepository().Create(&postCategory)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}
