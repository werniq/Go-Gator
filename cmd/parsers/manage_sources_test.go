package parsers

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	StoragePath = filepath.Join("..", "parsers", "data")
}

func TestAddNewSource(t *testing.T) {
	tests := []struct {
		format             string
		source             string
		endpoint           string
		finish             func()
		expectedEndpoint   string
		expectedParserType string
	}{
		{
			format: "json",
			finish: func() {
				err := DeleteSource("NewSourceJSON")
				assert.Nil(t, err)
			},
			source:             "NewSourceJSON",
			endpoint:           "https://newsample.com/json",
			expectedEndpoint:   "https://newsample.com/json",
			expectedParserType: "JsonParser",
		},
		{
			format: "xml",
			finish: func() {
				err := DeleteSource("NewSourceXML")
				assert.Nil(t, err)
			},
			source:             "NewSourceXML",
			endpoint:           "https://newsample.com/xml",
			expectedEndpoint:   "https://newsample.com/xml",
			expectedParserType: "XMLParser",
		},
		{
			format: "html",
			finish: func() {
				err := DeleteSource("NewSourceHTML")
				assert.Nil(t, err)
			},
			source:             "NewSourceHTML",
			endpoint:           "https://newsample.com/html",
			expectedEndpoint:   "https://newsample.com/html",
			expectedParserType: "HtmlParser",
		},
	}

	for _, tt := range tests {
		err := AddNewSource(tt.format, tt.source, tt.endpoint)
		assert.Nil(t, err)

		if sourceToEndpoint[tt.source] != tt.expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", tt.expectedEndpoint, sourceToEndpoint[tt.source])
		}

		parser := sourceToParser[tt.source]
		switch tt.expectedParserType {
		case "JsonParser":
			if _, ok := parser.(JsonParser); !ok {
				t.Errorf("expected JsonParser, got %T", parser)
			}
		case "XMLParser":
			if _, ok := parser.(XMLParser); !ok {
				t.Errorf("expected XMLParser, got %T", parser)
			}
		case "HtmlParser":
			if _, ok := parser.(HtmlParser); !ok {
				t.Errorf("expected HtmlParser, got %T", parser)
			}
		}
	}
}

func TestGetAllSources(t *testing.T) {
	expected := map[string]string{
		WashingtonTimes: "https://www.washingtontimes.com/rss/headlines/news/world",
		ABC:             "https://abcnews.go.com/abcnews/internationalheadlines",
		BBC:             "https://feeds.bbci.co.uk/news/rss.xml",
		UsaToday:        "https://usatoday.com",
	}

	result := GetAllSources()
	for key, value := range expected {
		if result[key] != value {
			t.Errorf("for key %s, expected %s, got %s", key, value, result[key])
		}
	}
}

func TestUpdateSourceEndpoint(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		newEndpoint      string
		expectedEndpoint string
		expectedErr      bool
	}{
		{
			name:             "Successful update",
			source:           "WashingtonTimes",
			newEndpoint:      "https://newendpoint.com/rss",
			expectedEndpoint: "https://newendpoint.com/rss",
			expectedErr:      false,
		},
		{
			name:             "Try to update not-existent source",
			source:           "source-not-exists",
			newEndpoint:      "https://api.com/rss",
			expectedEndpoint: "https://api.com/rss",
			expectedErr:      true,
		},
		{
			name:             "Try to update not-existent source",
			source:           "source-not-exists",
			newEndpoint:      "https://api.com/rss",
			expectedEndpoint: "https://api.com/rss",
			expectedErr:      true,
		},
	}

	for _, tt := range tests {
		err := UpdateSourceEndpoint(tt.source, tt.newEndpoint)

		if tt.expectedErr {
			assert.NotNil(t, err)

			if sourceToEndpoint[tt.source] != tt.expectedEndpoint {
				t.Errorf("for source %s, expected endpoint %s, got %s", tt.source, tt.expectedEndpoint, sourceToEndpoint[tt.source])
			}
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestUpdateSourceFormat(t *testing.T) {
	tests := []struct {
		name               string
		source             string
		format             string
		expectedParserType string
		expectedErr        bool
	}{
		{
			name:               "Successful update",
			source:             "WashingtonTimes",
			format:             "xml",
			expectedParserType: "xml",
			expectedErr:        false,
		},
		{
			name:               "Try to update not-existent source",
			source:             "source-not-exists",
			format:             "https://api.com/rss",
			expectedParserType: "https://api.com/rss",
			expectedErr:        true,
		},
		{
			name:               "Try to update not-existent source",
			source:             "source-not-exists",
			format:             "https://api.com/rss",
			expectedParserType: "https://api.com/rss",
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateSourceFormat(tt.source, tt.format)

			if tt.expectedErr {
				assert.NotNil(t, err)

				if sourceToEndpoint[tt.source] != tt.expectedParserType {
					t.Errorf("for source %s, expected endpoint %s, got %s", tt.source, tt.expectedParserType, sourceToEndpoint[tt.source])
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDeleteSource(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectedErr bool
	}{
		{
			name:        "Successful update",
			source:      "WashingtonTimes",
			expectedErr: false,
		},
		{
			name:        "Try to delete not-existent source",
			source:      "source-not-exists",
			expectedErr: true,
		},
		{
			name:        "Try to delete not-existent source",
			source:      "source-not-exists",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteSource(tt.source)

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDetermineParser(t *testing.T) {
	tests := []struct {
		format   string
		source   string
		expected Parser
	}{
		{
			format: "json",
			source: "source1",
			expected: JsonParser{
				Source: "source1",
			},
		},
		{
			format: "xml",
			source: "source2",
			expected: XMLParser{
				Source: "source2",
			},
		},
		{
			format: "html",
			source: "source3",
			expected: HtmlParser{
				Source: "source3",
			},
		},
		{
			format:   "invalid",
			source:   "source4",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := determineParser(tt.format, tt.source)
			if got != tt.expected {
				t.Errorf("determineParser(%s, %s) = %v, expected %v", tt.format, tt.source, got, tt.expected)
			}
		})
	}
}

func TestDetermineFormat(t *testing.T) {
	tests := []struct {
		parser   Parser
		source   string
		expected string
	}{
		{JsonParser{Source: "source1"}, "source1", "json"},
		{XMLParser{Source: "source2"}, "source2", "xml"},
		{HtmlParser{Source: "source3"}, "source3", "html"},
		{nil, "source4", ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := determineFormat(tt.parser, tt.source)
			if got != tt.expected {
				t.Errorf("determineFormat(%v, %s) = %v, expected %v", tt.parser, tt.source, got, tt.expected)
			}
		})
	}
}

func TestUpdateSourcesFile(t *testing.T) {
	tests := []struct {
		name          string
		setup         func()
		expectedError bool
	}{
		{
			name: "Successful execution",
			setup: func() {
			},
			expectedError: false,
		},
		{
			name: "Invalid (non-existent) path to the storage",
			setup: func() {
				StoragePath = "non/existent/path"
			},
			expectedError: true,
		},
		{
			name: "Remove data from sourceToParser, to cause error",
			setup: func() {
				sourceToParser = nil
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := UpdateSourcesFile()

			if tt.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
