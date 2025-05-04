package routes

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	router, cookies, _ := StartWithLogin()

	testsOK := map[string]map[string]int{
		"/admin/posts": {
			http.MethodGet: http.StatusOK,
		},
		"/admin": {
			http.MethodGet: http.StatusOK,
		},
		"/login": {
			http.MethodGet: http.StatusFound,
		},
	}

	for route, body := range testsOK {
		for method, status := range body {
			t.Run("Successful login", func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(method, route, nil)

				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}

				router.ServeHTTP(w, req)

				assert.Equal(t, status, w.Code)
			})
		}
	}

	tests := []string{
		``,
		`{"email":"axlle@mail","password":"123456"}`,
		`{"email":"axlle@mail.ru","password":"1"}`,
		`{"email":"","password":""}`,
		`{"email":"","password":"123456"}`,
		`{"emails":"axlle@mail.ru","passwords":"123456"}`,
	}

	for _, body := range tests {
		t.Run("Failed login", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusFound, w.Code)
			assert.Contains(t, w.Header().Get("Location"), "/login")
			assert.Contains(t, w.Body.String(), "")
		})
	}
}
