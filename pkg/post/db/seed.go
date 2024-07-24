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
	ids, _ := templateRepo.NewTemplateRepository().GetAllIds()
	idsCategory, _ := categoryRepo.NewPostCategoryRepository().GetAllIds()
	idsUser, _ := userRepo.NewUserRepository().GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomCategoryID := ids[rand.Intn(len(idsCategory))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]
		post := Post{
			TemplateID:         &randomID,
			PostCategoryID:     &randomCategoryID,
			UserID:             randomUserID,
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
			Image:              StrPtr(faker.Word()),
			Hits:               UintPtr(1000),
			Sort:               IntPtr(rand.Intn(100)),
			Stars:              Float32Ptr(rand.Float32() * 5),
			CreatedAt:          TimePtr(time.Now()),
			UpdatedAt:          TimePtr(time.Now()),
			DeletedAt:          nil,
		}

		err := NewRepository().CreatePost(&post)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Post successfully!")
}
