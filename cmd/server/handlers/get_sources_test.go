package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSources(t *testing.T) {
	server := gin.Default()
	server.GET("/admin/source", GetSources)

	tests := []struct {
		name       string
		source     string
		setup      func()
		statusCode int
		response   gin.H
	}{
		{
			name:       "Get all sources",
			statusCode: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/admin/source", nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
		})
	}
}
