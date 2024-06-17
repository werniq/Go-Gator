package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"strings"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		"usatoday",
	}
	params := &types.FilteringParams{Keywords: "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit"}
	news, err := parser.Parse()
	assert.Equal(t, err, nil)

	news = ApplyParams(news, params)

	assert.Len(t, news, 1, "Expected 1 news item")

	ok := strings.Contains(news[0].Description, "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit")
	assert.Equal(t, true, ok)

	ok = strings.Contains(news[0].Title, "The Taizhou Zoo in Jiangsu, China dyed two chow chow dogs and advertised them as \"panda dogs\" in the exhibit that opened on May 1..")
	assert.Equal(t, true, ok)
}
