package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoGatorParsingFactory_CreateHtmlParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	htmlParser := g.CreateHtmlParser()

	assert.Equal(t, htmlParser, HtmlParser{})
}

func TestGoGatorParsingFactory_CreateJsonParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	jsonParser := g.CreateJsonParser()

	assert.Equal(t, jsonParser, JsonParser{})
}

func TestGoGatorParsingFactory_CreateXmlParser(t *testing.T) {
	g := GoGatorParsingFactory{}

	xmlParser := g.CreateXmlParser()

	assert.Equal(t, xmlParser, XMLParser{})
}
