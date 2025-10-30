package routes

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service"
	modelsGallery "github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

var cntPost = 10
var cntGallery = 3
var cntImage = 10

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Post PostResponse `json:"post"`
}

type GalleryRequest struct {
	ID          string          `json:"id" form:"id" binding:"omitempty"`
	Title       string          `json:"title" form:"title" binding:"omitempty"`
	Description string          `json:"description" form:"description" binding:"omitempty"`
	Sort        string          `json:"sort" form:"sort" binding:"omitempty"`
	Image       string          `json:"image" form:"image" binding:"omitempty"`
	URL         string          `json:"url" form:"url" binding:"omitempty"`
	Images      []*ImageRequest `json:"images" form:"images" binding:"omitempty"`
}

type ImageRequest struct {
	ID           string `json:"id" form:"id" binding:"omitempty"`
	GalleryID    string `json:"gallery_id" form:"gallery_id" binding:"omitempty"`
	OriginalName string `json:"original_name" form:"original_name" binding:"omitempty"`
	File         string `json:"file" form:"file" binding:"omitempty"`
	Title        string `json:"title" form:"title" binding:"omitempty"`
	Description  string `json:"description" form:"description" binding:"omitempty"`
	Sort         string `json:"sort" form:"sort" binding:"omitempty"`
}

type PostRequest struct {
	ID                 string            `json:"id" form:"id"`
	TemplateID         string            `json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     string            `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          string            `json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    string            `json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string            `json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL                string            `json:"url" form:"url" binding:"omitempty,max=1000"`
	IsPublished        string            `json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       string            `json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        string            `json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      string            `json:"show_image_post" form:"show_image_post" binding:"omitempty"`
	ShowImageCategory  string            `json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	InSitemap          string            `json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              string            `json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string            `json:"title" form:"title" binding:"required,max=255"`
	TitleShort         string            `json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview string            `json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        string            `json:"description" form:"description" binding:"omitempty"`
	ShowDate           string            `json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            string            `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            string            `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              string            `json:"image" form:"image" binding:"omitempty,max=255"`
	Sort               string            `json:"sort" form:"sort" binding:"omitempty"`
	Galleries          []*GalleryRequest `json:"galleries" form:"galleries" binding:"omitempty"`
}

type PostResponse struct {
	ID                 uint                     `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	UserID             *uint                    `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         *uint                    `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     *uint                    `gorm:"index" json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          *string                  `gorm:"size:100" json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string                  `gorm:"size:200" json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string                   `gorm:"size:255;unique" json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL                string                   `gorm:"size:1000;unique" json:"url" form:"url" binding:"omitempty,max=1000"`
	IsPublished        bool                     `gorm:"not null;default:false" json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       bool                     `gorm:"not null;default:false" json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        bool                     `gorm:"not null;default:false" json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      bool                     `gorm:"not null;default:false" json:"show_image_post" form:"show_image_post"`
	ShowImageCategory  bool                     `gorm:"not null;default:false" json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	InSitemap          bool                     `gorm:"not null;default:false" json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              *string                  `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string                   `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	TitleShort         *string                  `gorm:"size:155;default:null" json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string                  `gorm:"type:text" json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        *string                  `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	ShowDate           bool                     `gorm:"not null;default:false" json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            *time.Time               `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            *time.Time               `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              *string                  `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	Hits               uint                     `gorm:"not null;default:0" json:"hits" form:"hits" binding:"-"`
	Sort               int                      `gorm:"not null;default:0" json:"sort" form:"sort" binding:"omitempty"`
	Stars              float32                  `gorm:"not null;default:0.0" json:"stars" form:"stars" binding:"-"`
	CreatedAt          *time.Time               `json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
	UpdatedAt          *time.Time               `json:"updated_at,omitempty" form:"updated_at" binding:"-" ignore:"true"`
	DeletedAt          *time.Time               `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-" ignore:"true"`
	Galleries          []*modelsGallery.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`
}

func TestFailedCreatePost(t *testing.T) {
	router, cookies, _ := StartWithLogin()
	requestBody := `{"email":"axlle@mail","password":"123456"}`

	t.Run("Failed login", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin/posts", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed create post", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/admin/posts", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"title is required"`)
		assert.Contains(t, w.Body.String(), `"message":"Ошибки валидации"`)
	})

	t.Run("Failed create post", func(t *testing.T) {
		requestBody = `{"title":""}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/admin/posts", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"title is required"`)
		assert.Contains(t, w.Body.String(), `"message":"Ошибки валидации"`)
	})
}

func TestSuccessfulCreatePost(t *testing.T) {
	router, cookies, _ := StartWithLogin()

	cfg := config.Config()
	newDB, _ := db.SetupDB(cfg)
	container := app.NewContainer(cfg, newDB)

	err := container.Migrator.Rollback()
	if err != nil {
		return
	}
	err = container.Migrator.Migrate()
	if err != nil {
		return
	}

	iProvider := container.ImageProvider
	gProvider := container.GalleryProvider
	pRepo := container.PostRepo

	for i := 0; i < cntPost; i++ {
		post := newPost()
		addNewGallery(post)

		t.Run("Successful create post", func(t *testing.T) {
			requestBody, err := json.Marshal(post)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/admin/posts", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var responseBody Response
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Error %v", err)
			}

			model, err := pRepo.GetByAlias(responseBody.Data.Post.Alias)
			if err != nil {
				t.Error(err)
			}

			gSlice := gProvider.GetForResourceUUID(model.UUID.String())
			if err != nil {
				t.Error(err)
			}
			if len(gSlice) != cntGallery {
				t.Errorf("Не верное количество постов %v", len(gSlice))
			}
			if len(post.Galleries) != cntGallery {
				t.Errorf("Не верное количество галлерей %v", len(post.Galleries))
			}

			postMap := make(map[string]*ImageRequest)
			modelMap := make(map[string]contract.Image)

			for _, gallery := range gSlice {
				iSlice := iProvider.GetForGallery(gallery.GetID())
				if len(iSlice) != cntImage {
					t.Errorf("Не верное количество изображений %v", len(iSlice))
				}

				for _, im := range iSlice {
					modelMap[im.GetFile()] = im
				}
			}

			for _, postGallery := range post.Galleries {
				for _, im := range postGallery.Images {
					if _, ok := modelMap[im.File]; ok {
						if im.Title != *modelMap[im.File].GetTitle() {
							t.Errorf("Не совпали Title ожидалось: %v; пришло: %v", im.Title, *modelMap[im.File].GetTitle())
						}

						if im.Description != *modelMap[im.File].GetDescription() {
							t.Errorf("Не совпали Description ожидалось: %v; пришло: %v", im.Description, *modelMap[im.File].GetDescription())
						}

						s, _ := service.ConvertStringToType(im.Sort, modelMap[im.File].GetSort())
						if s != modelMap[im.File].GetSort() {
							t.Errorf("Не совпали Sort ожидалось: %v; пришло: %v", im.Description, *modelMap[im.File].GetDescription())
						}

						delete(modelMap, im.File)
					} else {
						postMap[im.File] = im
					}
				}
			}

			if len(postMap) != len(modelMap) {
				t.Errorf("Не верное количество изображений %v", len(postMap))
			}

			var v any

			v, _ = service.ConvertStringToType(post.TemplateID, responseBody.Data.Post.TemplateID)
			assert.Equal(t, v, responseBody.Data.Post.TemplateID)
			assert.Equal(t, v, model.TemplateID)

			v, _ = service.ConvertStringToType(post.PostCategoryID, responseBody.Data.Post.PostCategoryID)
			assert.Equal(t, v, responseBody.Data.Post.PostCategoryID)
			assert.Equal(t, v, model.PostCategoryID)

			assert.Equal(t, post.Title, responseBody.Data.Post.Title)
			assert.Equal(t, post.Title, model.Title)

			v, _ = service.ConvertStringToType(post.MetaDescription, responseBody.Data.Post.MetaDescription)
			assert.Equal(t, v, responseBody.Data.Post.MetaDescription)
			assert.Equal(t, v, model.MetaDescription)

			v, _ = service.ConvertStringToType(post.IsPublished, responseBody.Data.Post.IsPublished)
			assert.Equal(t, v, responseBody.Data.Post.IsPublished)
			assert.Equal(t, v, model.IsPublished)

			v, _ = service.ConvertStringToType(post.IsFavourites, responseBody.Data.Post.IsFavourites)
			assert.Equal(t, v, responseBody.Data.Post.IsFavourites)
			assert.Equal(t, v, model.IsFavourites)

			v, _ = service.ConvertStringToType(post.HasComments, responseBody.Data.Post.HasComments)
			assert.Equal(t, v, responseBody.Data.Post.HasComments)
			assert.Equal(t, v, model.HasComments)

			v, _ = service.ConvertStringToType(post.ShowImagePost, responseBody.Data.Post.ShowImagePost)
			assert.Equal(t, v, responseBody.Data.Post.ShowImagePost)
			assert.Equal(t, v, model.ShowImagePost)

			v, _ = service.ConvertStringToType(post.ShowImageCategory, responseBody.Data.Post.ShowImageCategory)
			assert.Equal(t, v, responseBody.Data.Post.ShowImageCategory)
			assert.Equal(t, v, model.ShowImageCategory)

			v, _ = service.ConvertStringToType(post.InSitemap, responseBody.Data.Post.InSitemap)
			assert.Equal(t, v, responseBody.Data.Post.InSitemap)
			assert.Equal(t, v, model.InSitemap)

			v, _ = service.ConvertStringToType(post.Media, responseBody.Data.Post.Media)
			assert.Equal(t, v, responseBody.Data.Post.Media)
			assert.Equal(t, v, model.Media)

			v, _ = service.ConvertStringToType(post.Title, responseBody.Data.Post.Title)
			assert.Equal(t, v, responseBody.Data.Post.Title)
			assert.Equal(t, v, model.Title)

			v, _ = service.ConvertStringToType(post.TitleShort, responseBody.Data.Post.TitleShort)
			assert.Equal(t, v, responseBody.Data.Post.TitleShort)
			assert.Equal(t, v, model.TitleShort)

			v, _ = service.ConvertStringToType(post.Description, responseBody.Data.Post.Description)
			assert.Equal(t, v, responseBody.Data.Post.Description)
			assert.Equal(t, v, model.Description)

			v, _ = service.ConvertStringToType(post.DescriptionPreview, responseBody.Data.Post.DescriptionPreview)
			assert.Equal(t, v, responseBody.Data.Post.DescriptionPreview)
			assert.Equal(t, v, model.DescriptionPreview)

			v, _ = service.ConvertStringToType(post.ShowDate, responseBody.Data.Post.ShowDate)
			assert.Equal(t, v, responseBody.Data.Post.ShowDate)
			assert.Equal(t, v, model.ShowDate)

			assert.Equal(t, post.DatePub, db.FormatDate(*responseBody.Data.Post.DatePub))
			assert.Equal(t, post.DatePub, db.FormatDate(*model.DatePub))

			assert.Equal(t, post.DateEnd, db.FormatDate(*responseBody.Data.Post.DateEnd))
			assert.Equal(t, post.DateEnd, db.FormatDate(*model.DateEnd))
		})
	}

	modelSlice, err := pRepo.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(modelSlice) != cntPost {
		t.Errorf("Не верное количество постов %v", len(modelSlice))
	}

	gSlice := gProvider.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(gSlice) != cntPost*cntGallery {
		t.Errorf("Не верное количество галлерей %v", len(gSlice))
	}

	iSlice := iProvider.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(iSlice) != cntPost*cntGallery*cntImage {
		t.Errorf("Не верное количество изображений %v", len(iSlice))
	}
}

func newPost() *PostRequest {
	post := &PostRequest{
		TemplateID:         db.IntStr(rand.Intn(10)),
		PostCategoryID:     db.IntStr(rand.Intn(10)),
		MetaTitle:          faker.Sentence(),
		MetaDescription:    faker.Sentence(),
		IsPublished:        db.IntStr(rand.Intn(2)),
		IsFavourites:       db.IntStr(rand.Intn(2)),
		HasComments:        db.IntStr(rand.Intn(2)),
		ShowImagePost:      db.IntStr(rand.Intn(2)),
		ShowImageCategory:  db.IntStr(rand.Intn(2)),
		InSitemap:          db.IntStr(rand.Intn(2)),
		Media:              faker.Sentence(),
		Title:              faker.Sentence(),
		TitleShort:         faker.Sentence(),
		DescriptionPreview: faker.Sentence(),
		Description:        faker.Sentence(),
		ShowDate:           db.IntStr(rand.Intn(2)),
		DatePub:            db.FormatDate(db.RandomDate()),
		DateEnd:            db.FormatDate(db.RandomDate()),
		Image:              faker.Sentence(),
		Sort:               db.IntStr(rand.Intn(20)),
	}
	return post
}

func addNewGallery(p *PostRequest) {
	for i := 0; i < cntGallery; i++ {
		g := gallery()
		p.Galleries = append(p.Galleries, g)
	}
}

func updateGallery(p *PostRequest, new []*modelsGallery.Gallery) {
	var galleries []*GalleryRequest
	for _, g := range new {
		newG := gallery()
		newG.ID = strconv.Itoa(int(g.ID))
		for i, im := range g.Images {
			newG.Images[i].ID = strconv.Itoa(int(im.ID))
		}
		galleries = append(galleries, newG)
	}
	p.Galleries = galleries
}

func gallery() *GalleryRequest {
	gallery := &GalleryRequest{
		Title:       faker.Sentence(),
		Description: faker.Sentence(),
		Image:       faker.Sentence(),
		Sort:        db.IntStr(rand.Intn(20)),
	}
	addNewImage(gallery)
	return gallery
}

func addNewImage(g *GalleryRequest) {
	var images []*ImageRequest
	for i := 0; i < cntImage; i++ {
		im := image()
		images = append(images, im)
	}
	g.Images = images
}

func image() *ImageRequest {
	return &ImageRequest{
		OriginalName: faker.Sentence(),
		File:         faker.UUIDHyphenated(),
		Title:        faker.Sentence(),
		Description:  faker.Sentence(),
		Sort:         db.IntStr(rand.Intn(20)),
	}
}
