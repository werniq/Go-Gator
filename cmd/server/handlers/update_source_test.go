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
	// Initialize Gin engine
	server := gin.Default()
	server.PUT("/admin/source", UpdateSource)

	tests := []struct {
		name       string
		body       types.Source
		source     string
		setup      func()
		statusCode int
		response   gin.H
	}{
		{
			name:   "Update non-existent source",
			source: "source6",
			body: types.Source{
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
			name:   "Update format in existent source",
			source: "source5",
			body: types.Source{
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
			body: types.Source{
				Name:     "bbc",
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
			reqBody, _ = json.Marshal(tt.body)

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
