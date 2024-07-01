package parsers

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAddNewSource(t *testing.T) {
	tests := []struct {
		format             string
		source             string
		endpoint           string
		expectedEndpoint   string
		expectedParserType string
	}{
		{"json", "NewSourceJSON", "https://newsample.com/json", "https://newsample.com/json", "JsonParser"},
		{"xml", "NewSourceXML", "https://newsample.com/xml", "https://newsample.com/xml", "XMLParser"},
		{"html", "NewSourceHTML", "https://newsample.com/html", "https://newsample.com/html", "HtmlParser"},
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
		NbcNews:         "nbc-news.json",
	}

	result := GetAllSources()
	for key, value := range expected {
		log.Println(result)
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
		UpdateSourceEndpoint(tt.source, tt.newEndpoint)

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
		UpdateSourceFormat(tt.source, tt.format)

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
