package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteSource(t *testing.T) {
	// Initialize Gin engine
	server := gin.Default()
	server.DELETE("/admin/source", DeleteSource)

	tests := []struct {
		name       string
		source     string
		setup      func()
		statusCode int
		response   gin.H
	}{
		{
			name:   "Delete existing source",
			source: "source1",
			setup: func() {
				err := parsers.AddNewSource("xml", "source1", "https://source1.com")
				if err != nil {
					assert.Equal(t, err, nil)
				}
			},
			statusCode: http.StatusOK,
			response: gin.H{
				"status": MsgSourceDeleted,
			},
		},
		{
			name:   "Delete non-existent source",
			source: "source2",
			setup: func() {
				err := parsers.DeleteSource("source2")
				if err != nil {
					assert.Equal(t, err, nil)
				}
			},
			statusCode: http.StatusBadRequest,
			response: gin.H{
				"error": ErrSourceNotFound,
			},
		},
		{
			name:       "Invalid JSON",
			source:     "",
			setup:      func() {},
			statusCode: http.StatusInternalServerError,
			response: gin.H{
				"error": ErrFailedToDecode + "invalid character 'i' looking for beginning of object key string",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()

			var reqBody []byte
			if testCase.source != "" {
				reqBody, _ = json.Marshal(gin.H{"name": testCase.source})
			} else {
				reqBody = []byte("{invalid json")
			}

			req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/admin/source", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)

			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.response, response)
		})
	}
}
