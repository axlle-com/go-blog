package testutils

import (
	"bytes"
	"encoding/json"
	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	"github.com/axlle-com/blog/pkg/gallery/provider"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	modelPost "github.com/axlle-com/blog/pkg/post/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFailedDeletePost(t *testing.T) {
	router := SetupTestRouter()
	cookies, _ := PerformLogin(router)

	t.Run("Failed delete post", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/admin/posts/200000000", nil)
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"message":"Ресурс не найден"`)
	})

	t.Run("Failed delete post", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/admin/posts/200000000", nil)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSuccessfulDeletePost(t *testing.T) {
	cookies, _ := PerformLogin(router)
	router := SetupTestRouter()
	mGallery.Rollback()
	mGallery.Migrate()
	mPost.Rollback()
	mPost.Migrate()

	modelSlice, err := modelPost.PostRepo().GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(modelSlice) != 0 {
		t.Errorf("Не верное количество постов %v", len(modelSlice))
	}

	gSlice := provider.Provider().GetAll()
	if len(gSlice) != 0 {
		t.Errorf("Не верное количество галлерей ожидалось: %v; итого: %v", 0, len(gSlice))
	}

	iSlice := provider.ImageProvider().GetAll()
	if len(iSlice) != 0 {
		t.Errorf("Не верное количество изображений ожидалось: %v; итого: %v", 0, len(iSlice))
	}

	for i := 0; i < cntPost; i++ {
		post := newPost()
		addNewGallery(post)

		var oPost *PostResponse

		t.Run("Successful delete post", func(t *testing.T) {
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

		t.Run("Successful delete post", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/admin/posts/"+strconv.Itoa(int(oPost.ID)), nil)
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}

	modelSlice, err = modelPost.PostRepo().GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(modelSlice) != 0 {
		t.Errorf("Не верное количество постов %v", len(modelSlice))
	}

	gSlice = provider.Provider().GetAll()
	if len(gSlice) != 0 {
		t.Errorf("Не верное количество галлерей ожидалось: %v; итого: %v", 0, len(gSlice))
	}

	iSlice = provider.ImageProvider().GetAll()
	if len(iSlice) != 0 {
		t.Errorf("Не верное количество изображений ожидалось: %v; итого: %v", 0, len(iSlice))
	}
}