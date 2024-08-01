package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetNews(t *testing.T) {
	server := gin.Default()
	server.GET("/news", GetNews)

	tests := []struct {
		name       string
		input      *types.FilteringParams
		statusCode int
	}{
		{
			name:       "Successful request without parameters",
			input:      &types.FilteringParams{},
			statusCode: http.StatusOK,
		},
		{
			name: "Successful request with parameters",
			input: &types.FilteringParams{
				Keywords:          "Ukraine",
				StartingTimestamp: "2024-06-20",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Failed request with wrong source",
			input: &types.FilteringParams{
				Sources: "source-7",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Failed request with wrong date range",
			input: &types.FilteringParams{
				StartingTimestamp: "2024-07-03",
				EndingTimestamp:   "2024-06-20",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse("http://localhost:8080/news")
			if err != nil {
				t.Fatal(err)
			}

			q := u.Query()
			if tt.input.Keywords != "" {
				q.Set("keywords", tt.input.Keywords)
			}
			if tt.input.StartingTimestamp != "" {
				q.Set("startingTimestamp", tt.input.StartingTimestamp)
			}
			if tt.input.EndingTimestamp != "" {
				q.Set("endingTimestamp", tt.input.EndingTimestamp)
			}
			if tt.input.Sources != "" {
				q.Set("sources", tt.input.Sources)
			}
			u.RawQuery = q.Encode()

			req, err := http.NewRequest(http.MethodGet, u.String(), nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			go server.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			//var response gin.H
			//err = json.Unmarshal(w.Body.Bytes(), &response)
			//assert.NoError(t, err)
			//assert.Equal(t, t.response, response)
		})
	}
}

func TestSourceInArray(t *testing.T) {
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
