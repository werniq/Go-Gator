package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSourceDetailed(t *testing.T) {
	server := gin.Default()
	server.GET("/admin/source/", GetSourceDetailed)
	server.GET("/admin/source/:source", GetSourceDetailed)

	tests := []struct {
		name       string
		source     string
		statusCode int
	}{
		{
			name:       "Get detailed information about bbc",
			source:     "bbc",
			statusCode: http.StatusOK,
		},
		{
			name:       "Try to get wrong source",
			source:     "no-source",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Try to get source with no name",
			source:     "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/admin/source/"+tt.source, nil)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
