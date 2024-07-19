package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedPosts(n int) {
	for i := 0; i < n; i++ {
		post := Post{
			TemplateID:         UintPtr(100),
			PostCategoryID:     UintPtr(100),
			MetaTitle:          StrPtr(faker.Sentence()),
			MetaDescription:    StrPtr(faker.Sentence()),
			Alias:              faker.Username(),
			URL:                faker.URL(),
			IsPublished:        IntToBoolPtr(),
			IsFavourites:       IntToBoolPtr(),
			IsComments:         IntToBoolPtr(),
			IsImagePost:        IntToBoolPtr(),
			IsImageCategory:    IntToBoolPtr(),
			IsWatermark:        IntToBoolPtr(),
			IsSitemap:          IntToBoolPtr(),
			Media:              StrPtr(faker.Word()),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			PreviewDescription: StrPtr(faker.Paragraph()),
			Description:        StrPtr(faker.Paragraph()),
			ShowDate:           IntToBoolPtr(),
			DatePub:            ParseDate(faker.Date()),
			DateEnd:            ParseDate(faker.Date()),
			ControlDatePub:     IntToBoolPtr(),
			ControlDateEnd:     IntToBoolPtr(),
			Image:              StrPtr(faker.Word()),
			Hits:               UintPtr(1000),
			Sort:               IntPtr(rand.Intn(100)),
			Stars:              Float32Ptr(rand.Float32() * 5),
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := NewPostRepository().CreatePost(&post)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}
