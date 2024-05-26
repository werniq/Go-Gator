package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
	"strings"
	"testing"
)

func TestXMLParser_Parse(t *testing.T) {
	params := &types.FilteringParams{
		Keywords: "Ukraine and Russia exchange drone attacks while Russia continues its push in the east",
	}
	parser := XMLParser{}
	news := parser.Parse(params)

	assert.Len(t, news, 1, "Expected 1 news item")

	ok := strings.Contains(news[0].Title, "Ukraine and Russia exchange drone attacks while Russia continues its push in the east")
	assert.Equal(t, true, ok)

	ok = strings.Contains(news[0].Description, "At least 10 people were reported killed in attacks in Ukraine&rsquo;s war-ravaged northeast on Sunday, as Russia pushed ahead with its renewed offensive")
	assert.Equal(t, true, ok)
}
