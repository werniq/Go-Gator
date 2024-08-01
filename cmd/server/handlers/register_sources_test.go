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

func TestRegisterSource(t *testing.T) {
	// Initialize Gin engine
	server := gin.Default()
	server.POST("/admin/source", RegisterSource)

	tests := []struct {
		name       string
		source     string
		setup      func()
		finish     func()
		statusCode int
		response   gin.H
	}{
		{
			name:   "Register non-existent source",
			source: "source5",
			setup:  func() {},
			finish: func() {
				err := parsers.DeleteSource("source5")
				assert.Nil(t, err)
			},
			statusCode: http.StatusCreated,
			response: gin.H{
				"status": MsgSourceCreated,
			},
		},
		{
			name:   "Register existent source",
			source: "source3",
			setup: func() {
				err := parsers.AddNewSource("xml", "source3", "https://source1.com")
				if err != nil {
					assert.Nil(t, err)
				}
			},
			finish: func() {
				err := parsers.DeleteSource("source3")
				assert.Nil(t, err)
			},
			statusCode: http.StatusBadRequest,
			response: gin.H{
				"error": ErrSourceExists,
			},
		},
		{
			name:       "Invalid JSON",
			source:     "",
			setup:      func() {},
			finish:     func() {},
			statusCode: http.StatusInternalServerError,
			response: gin.H{
				"error": ErrFailedToDecode + "invalid character 't' looking for beginning of object key string",
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
				reqBody = []byte("{testing with invalid json")
			}

			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/admin/source", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)

			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.response, response)

			testCase.finish()
		})
	}
}
