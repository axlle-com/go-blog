package testutils

import (
	"bytes"
	"encoding/json"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/common/service"
	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFailedUpdatePost(t *testing.T) {
	router, cookies, _ := StartWithLogin()
	requestBody := `{"title":"title"}`
	var oPost *PostResponse

	t.Run("Failed update post", func(t *testing.T) {
		requestBody = `{"title":"title"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/admin/posts/200000000", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"message":"Ресурс не найден"`)
	})

	t.Run("Successful create post", func(t *testing.T) {
		post := newPost()
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

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody Response
		err = json.Unmarshal(w.Body.Bytes(), &responseBody)
		if err != nil {
			t.Errorf("Error %v", err)
		}

		oPost = &responseBody.Data.Post
	})

	t.Run("Failed update post", func(t *testing.T) {
		requestBody = `{"title":""}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/admin/posts/"+strconv.Itoa(int(oPost.ID)), bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"title is required"`)
		assert.Contains(t, w.Body.String(), `"message":"Ошибки валидации"`)
	})

	t.Run("Failed update post", func(t *testing.T) {
		text := faker.Paragraph()
		requestBody = `{"title":"` + text + text + text + text + text + text + text + text + `"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/admin/posts/"+strconv.Itoa(int(oPost.ID)), bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"title must be at most 255 characters long"`)
		assert.Contains(t, w.Body.String(), `"message":"Ошибки валидации"`)
	})
}

func TestSuccessfulUpdatePost(t *testing.T) {
	router, cookies, _ := StartWithLogin()
	mGallery.Rollback()
	mGallery.Migrate()
	mPost.Rollback()
	mPost.Migrate()

	for i := 0; i < cntPost; i++ {
		post := newPost()
		addNewGallery(post)

		var oPost *PostResponse

		t.Run("Successful update post", func(t *testing.T) {
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

			assert.Equal(t, http.StatusOK, w.Code)

			var responseBody Response
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Error %v", err)
			}

			oPost = &responseBody.Data.Post
		})

		t.Run("Successful update post", func(t *testing.T) {
			post := newPost()
			updateGallery(post, oPost.Galleries)
			addNewGallery(post)

			post.ID = strconv.Itoa(int(oPost.ID))
			requestBody, err := json.Marshal(post)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/admin/posts/"+post.ID, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var responseBody Response
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Error %v", err)
			}

			model, err := models.PostRepo().GetByAlias(responseBody.Data.Post.Alias)
			if err != nil {
				t.Error(err)
			}

			gSlice := provider.Provider().GetForResource(model)
			if err != nil {
				t.Error(err)
			}
			if len(gSlice) != cntGallery*2 {
				t.Errorf("Не верное количество галлерей %v", len(gSlice))
			}
			if len(post.Galleries) != cntGallery*2 {
				t.Errorf("Не верное количество галлерей %v", len(post.Galleries))
			}

			postMap := make(map[string]*ImageRequest)
			modelMap := make(map[string]contracts.Image)

			for _, gallery := range gSlice {
				iSlice := provider.ImageProvider().GetForGallery(gallery.GetID())
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

	gSlice := provider.Provider().GetAll()
	if len(gSlice) != cntPost*cntGallery*2 {
		t.Errorf("Не верное количество галлерей ожидалось: %v; итого: %v", cntPost*cntGallery*2, len(gSlice))
	}

	iSlice := provider.ImageProvider().GetAll()
	if len(iSlice) != cntPost*cntGallery*cntImage*2 {
		t.Errorf("Не верное количество изображений ожидалось: %v; итого: %v", cntPost*cntGallery*cntImage*2, len(iSlice))
	}
}

func TestSuccessfulUpdatePostAlias(t *testing.T) {
	router, cookies, _ := StartWithLogin()
	mPost.Rollback()
	mPost.Migrate()

	sliceCreate := []map[string]string{
		{
			"Alias":         "consequatur-aut-sit-voluptatem-accusantium-perferendis",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
		{
			"Alias":         "consequatur-aut-sit-voluptatem-accusantium-perferendis",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis-1",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
		{
			"Alias":         "",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis-2",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
		{
			"Alias":         "consequatur-aut-sit-voluptatem-accusantium-perferendis",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis-3",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
		{
			"Alias":         "",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis-4",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
		{
			"Alias":         "",
			"Title":         "Consequatur aut sit voluptatem accusantium perferendis.",
			"AliasExpected": "consequatur-aut-sit-voluptatem-accusantium-perferendis-5",
			"TitleExpected": "Consequatur aut sit voluptatem accusantium perferendis.",
		},
	}

	for i := 0; i < len(sliceCreate); i++ {
		post := newPost()

		var oPost *PostResponse

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

			assert.Equal(t, http.StatusOK, w.Code)

			var responseBody Response
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Error %v", err)
			}

			oPost = &responseBody.Data.Post
		})

		t.Run("Successful update post", func(t *testing.T) {
			post.ID = strconv.Itoa(int(oPost.ID))
			post.Alias = sliceCreate[i]["Alias"]
			post.Title = sliceCreate[i]["Title"]

			requestBody, err := json.Marshal(post)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/admin/posts/"+post.ID, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var responseBody Response
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Error %v", err)
			}

			model, err := models.PostRepo().GetByAlias(responseBody.Data.Post.Alias)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, sliceCreate[i]["AliasExpected"], responseBody.Data.Post.Alias)
			assert.Equal(t, sliceCreate[i]["TitleExpected"], responseBody.Data.Post.Title)

			var v any

			v, _ = service.ConvertStringToType(post.TemplateID, responseBody.Data.Post.TemplateID)
			assert.Equal(t, v, responseBody.Data.Post.TemplateID)
			assert.Equal(t, v, model.TemplateID)

			v, _ = service.ConvertStringToType(post.PostCategoryID, responseBody.Data.Post.PostCategoryID)
			assert.Equal(t, v, responseBody.Data.Post.PostCategoryID)
			assert.Equal(t, v, model.PostCategoryID)

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
	mGallery.Rollback()
	mPost.Rollback()
}
