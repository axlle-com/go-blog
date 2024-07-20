package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/repository"
	"github.com/axlle-com/blog/pkg/search"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedPosts(n int) {
	client := search.NewElasticsearch()
	err := client.CreateIndex("post")
	if err != nil {
		log.Printf("Failed to create index: %v", err.Error())
	}
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
			HasComments:        IntToBoolPtr(),
			ShowImagePost:      IntToBoolPtr(),
			ShowImageCategory:  IntToBoolPtr(),
			MakeWatermark:      IntToBoolPtr(),
			InSitemap:          IntToBoolPtr(),
			Media:              StrPtr(faker.Word()),
			Title:              faker.Sentence(),
			TitleShort:         StrPtr(faker.Sentence()),
			DescriptionPreview: StrPtr(faker.Paragraph()),
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
		} else {
			err := client.AddPost(&post)
			if err != nil {
				continue
			}
		}
	}

	log.Println("Database seeded Post successfully!")
}
