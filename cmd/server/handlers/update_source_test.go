package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"newsaggr/cmd/types"
	"testing"
)

func TestUpdateSource(t *testing.T) {
	server := gin.Default()
	server.PUT("/admin/source", UpdateSource)

	tests := []struct {
		name       string
		body       *types.Source
		source     string
		setup      func()
		statusCode int
		response   gin.H
	}{
		{
			name:   "Update non-existent source",
			source: "source6",
			body: &types.Source{
				Name:     "source6",
				Format:   "xml",
				Endpoint: "https://source5.com",
			},
			setup:      func() {},
			statusCode: http.StatusBadRequest,
			response: gin.H{
				"error": ErrSourceNotFound,
			},
		},
		{
			name:   "Empty source name",
			source: "source6",
			body: &types.Source{
				Name:     "",
				Format:   "",
				Endpoint: "",
			},
			setup:      func() {},
			statusCode: http.StatusBadRequest,
			response: gin.H{
				"error": ErrNoSourceName,
			},
		},
		{
			name:       "Empty source struct",
			source:     "source6",
			body:       nil,
			setup:      func() {},
			statusCode: http.StatusBadRequest,
			response: gin.H{
				"error": ErrFailedToDecode + "json: cannot unmarshal string into Go value of type types.Source",
			},
		},
		{
			name:   "Update format in existent source",
			source: "source5",
			body: &types.Source{
				Name:   "bbc",
				Format: "html",
			},
			setup:      func() {},
			statusCode: http.StatusOK,
			response: gin.H{
				"status": MsgSourceUpdated,
			},
		},
		{
			name:   "Update endpoint in existent source",
			source: "source5",
			body: &types.Source{
				Name:     "bbc",
				Format:   "",
				Endpoint: "https://bbc.com/",
			},
			setup:      func() {},
			statusCode: http.StatusOK,
			response: gin.H{
				"status": MsgSourceUpdated,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			} else {
				reqBody, _ = json.Marshal("{invalid json")
			}

			req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/admin/source", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.response, response)
		})
	}
}
