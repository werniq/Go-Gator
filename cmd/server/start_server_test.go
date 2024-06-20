package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRunServer(t *testing.T) {
	// This test ensures that ConfAndRun does not return an error on startup
	go func() {
		ConfAndRun()
	}()

	// Issue a request to ensure the server is running
	resp, err := http.Get("http://localhost:8080/news?sources=abc")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
