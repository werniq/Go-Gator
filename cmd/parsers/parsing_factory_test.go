package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoGatorParsingFactory_CreateHtmlParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	htmlParser := g.CreateHtmlParser("source")

	assert.Equal(t, htmlParser, HtmlParser{"source"})
}

func TestGoGatorParsingFactory_CreateJsonParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	jsonParser := g.CreateJsonParser("source")

	assert.Equal(t, jsonParser, JsonParser{"source"})
}

func TestGoGatorParsingFactory_CreateXmlParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	xmlParser := g.CreateXmlParser("source")

	assert.Equal(t, xmlParser, XMLParser{"source"})
}
