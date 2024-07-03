package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"newsaggr/cmd/parsers"
	"testing"
)

func TestSourceInArray(t *testing.T) {
	// Test cases
	tests := []struct {
		source   string
		expected bool
	}{
		{parsers.WashingtonTimes, true},
		{parsers.ABC, true},
		{parsers.BBC, true},
		{parsers.UsaToday, true},
		{"no-source", false},
		{"abcbbccbb", false},
	}

	for _, test := range tests {
		result := sourceInArray(test.source)
		assert.Equal(t, test.expected, result, "Expected %v for source %s, got %v", test.expected, test.source, result)
	}
}

func TestGetNews(t *testing.T) {
	// Initialize Gin engine
	r := gin.Default()
	r.GET("/admin/source", GetSources)

	tests := []struct {
		name       string
		statusCode int
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
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
		})
	}
}
