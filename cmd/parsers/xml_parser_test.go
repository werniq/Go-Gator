package parsers

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/types"
	"testing"
)

func TestXmlParser_Parse(t *testing.T) {
	parser := XMLParser{
		Source: "abc",
	}

	testCases := []struct {
		name         string
		setupMock    func()
		expectError  bool
		expectedNews []types.News
	}{
		{
			name: "Default parse",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				mockXML := `<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0">
					<channel>
						<title>Test Channel</title>
						<item>
							<title>Test News</title>
							<description>This is a test news.</description>
							<pubDate>2024-07-23</pubDate>
							<link>http://example.com</link>
						</item>
					</channel>
				</rss>`

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(200, mockXML))
			},
			expectError: false,
			expectedNews: []types.News{
				{
					Title:       "Test News",
					Description: "This is a test news.",
					PubDate:     "2024-07-23",
					Publisher:   "abc",
					Link:        "http://example.com",
				},
			},
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
					httpmock.NewStringResponder(200, ""))
			},
			expectError: true,
		},
		{
			name: "Invalid XML format",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(
						200,
						`<rss><channel><item><title>Test</title></item></channel></rss>`))
			},
			expectError: true,
		},
		{
			name: "Unexpected XML structure",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(
						200,
						`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel><title>Test Channel</title></channel></rss>`))
			},
			expectError:  false,
			expectedNews: []types.News{},
		},
		{
			name: "Empty XML document",
			setupMock: func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", sourceToEndpoint[parser.Source],
					httpmock.NewStringResponder(
						200,
						`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel></channel></rss>`))
			},
			expectError:  false,
			expectedNews: []types.News{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			news, err := parser.Parse()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedNews, news)
			}
		})
	}
}
