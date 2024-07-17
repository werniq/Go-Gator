package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoGatorParsingFactory_CreateHtmlParser(t *testing.T) {
	g := ParsingFactory{}

	htmlParser := g.HtmlParser("source")

	assert.Equal(t, htmlParser, HtmlParser{"source"})
}

func TestGoGatorParsingFactory_CreateJsonParser(t *testing.T) {
	g := ParsingFactory{}

	jsonParser := g.JsonParser("source")

	assert.Equal(t, jsonParser, JsonParser{"source"})
}

func TestGoGatorParsingFactory_CreateXmlParser(t *testing.T) {
	g := ParsingFactory{}

	xmlParser := g.XmlParser("source")

	assert.Equal(t, xmlParser, XMLParser{"source"})
}
