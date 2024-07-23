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
			expectedParserType: "XMLParser"},
		{
			format: "html",
			finish: func() {
				err := DeleteSource("NewSourceHTML")
				assert.Nil(t, err)
			},
			source:             "NewSourceHTML",
			endpoint:           "https://newsample.com/html",
			expectedEndpoint:   "https://newsample.com/html",
			expectedParserType: "HtmlParser"},
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
		source           string
		newEndpoint      string
		expectedEndpoint string
	}{
		{"WashingtonTimes", "https://newendpoint.com/rss", "https://newendpoint.com/rss"},
		{"ABC", "https://newendpoint.com/abc", "https://newendpoint.com/abc"},
	}

	for _, tt := range tests {
		err := UpdateSourceEndpoint(tt.source, tt.newEndpoint)
		assert.Nil(t, err)

		if sourceToEndpoint[tt.source] != tt.expectedEndpoint {
			t.Errorf("for source %s, expected endpoint %s, got %s", tt.source, tt.expectedEndpoint, sourceToEndpoint[tt.source])
		}
	}
}

func TestUpdateSourceFormat(t *testing.T) {
	tests := []struct {
		source             string
		format             string
		expectedParserType string
	}{
		{"WashingtonTimes", "json", "JsonParser"},
		{"ABC", "html", "HtmlParser"},
	}

	for _, tt := range tests {
		err := UpdateSourceFormat(tt.source, tt.format)
		assert.Nil(t, err)

		parser := sourceToParser[tt.source]
		switch tt.expectedParserType {
		case "JsonParser":
			if _, ok := parser.(JsonParser); !ok {
				t.Errorf("expected JsonParser, got %T", parser)
			}
		case "HtmlParser":
			if _, ok := parser.(HtmlParser); !ok {
				t.Errorf("expected HtmlParser, got %T", parser)
			}
		}
	}
}

func TestDetermineParser(t *testing.T) {
	tests := []struct {
		format   string
		source   string
		expected Parser
	}{
		{"json", "source1", JsonParser{Source: "source1"}},
		{"xml", "source2", XMLParser{Source: "source2"}},
		{"html", "source3", HtmlParser{Source: "source3"}},
		{"invalid", "source4", nil},
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
