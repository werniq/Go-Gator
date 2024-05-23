package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
	"strings"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	params := &types.ParsingParams{
		Keywords: "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit",
	}
	parser := HtmlParser{}
	news := parser.Parse(params)

	assert.Len(t, news, 1, "Expected 1 news item")

	ok := strings.Contains(news[0].Description, "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit")
	assert.Equal(t, true, ok)

	ok = strings.Contains(news[0].Title, "The Taizhou Zoo in Jiangsu, China dyed two chow chow dogs and advertised them as \"panda dogs\" in the exhibit that opened on May 1..")
	assert.Equal(t, true, ok)
}
