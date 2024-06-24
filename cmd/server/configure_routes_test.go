package server

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	router := gin.Default()
	setupRoutes(router)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "http://localhost:8080/news", http.StatusOK},
		{"GET", "/admin/source", http.StatusOK},
		{"POST", "/admin/source", http.StatusOK},
		{"PUT", "/admin/source", http.StatusOK},
		{"DELETE", "/admin/source", http.StatusOK},
		{"GET", "/nonexistent", http.StatusNotFound},
		{"POST", "/admin/source/extra", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.Nil(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
