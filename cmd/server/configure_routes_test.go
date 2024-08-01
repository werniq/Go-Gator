package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{"GET /news", "GET", "/non-existent", http.StatusNotFound},
		{"GET /admin/sources", "GET", "/admin/sources", http.StatusOK},
		{"PUT /admin/sources", "PUT", "/admin/sources", http.StatusOK},
		{"POST /admin/sources", "POST", "/admin/sources", http.StatusOK},
		{"DELETE /admin/sources", "DELETE", "/admin/sources", http.StatusOK},
		{"GET /news", "GET", "/news", http.StatusOK},
		{"POST /news", "GET", "/news", http.StatusNotFound},
		{"DELETE /news", "GET", "/news", http.StatusNotFound},
	}

	gin.SetMode(gin.TestMode)
	server := gin.Default()
	setupRoutes(server)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte{}))

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.statusCode, resp.Code)
		})
	}
}
