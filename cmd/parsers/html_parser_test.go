package parsers

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		Source: UsaToday,
	}

	testCases := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "Default run",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(
						http.StatusOK,
						`<!DOCTYPE html><html><head><title>Test</title></head><body><div class="news-item"><h1 class="title">Test News</h1><time datetime="2024-07-23">July 23, 2024</time><a href="/test-link">Link</a><p>Description</p></div></body></html>`))
			},
			expectError: false,
		},
		{
			name: "HTTP request failure",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewErrorResponder(errors.New("http request failed")))
			},
			expectError: true,
		},
		{
			name: "Empty response body",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(http.StatusOK, ""))
			},
			expectError: true,
		},
		{
			name: "Invalid HTML",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(http.StatusOK, "<html><head><title>Test</title></head><body><div><h1>Invalid HTML</h1></div>"))
			},
			expectError: true,
		},
		{
			name: "Missing selectors",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(http.StatusOK, `<!DOCTYPE html><html><head><title>Test</title></head><body><div><h1>No News Item</h1></div></body></html>`))
			},
			expectError: false,
		},
		{
			name: "Attributes missing",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(http.StatusOK, `<!DOCTYPE html><html><head><title>Test</title></head><body><div class="news-item"><h1 class="title">Test News</h1><p>Description</p></div></body></html>`))
			},
			expectError: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := parser.Parse()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
