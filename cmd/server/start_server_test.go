package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRunServer(t *testing.T) {
	testCases := []struct {
		Url          string
		ExpectedCode int
	}{
		{
			Url:          "http://localhost:8080",
			ExpectedCode: http.StatusOK,
		},
		{
			Url:          "http://localhost:8080/some-path",
			ExpectedCode: http.StatusNotFound,
		},
	}

	go func() {
		ConfAndRun()
	}()

	for _, tt := range testCases {
		resp, err := http.Get(tt.Url)
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, tt.ExpectedCode)
	}
}
