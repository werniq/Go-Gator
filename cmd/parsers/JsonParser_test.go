package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
	"testing"
)

func TestJSONParser_Parse(t *testing.T) {
	params := &types.ParsingParams{
		Keywords: "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena",
	}

	parser := JsonParser{}
	news := parser.Parse(params)

	assert.Len(t, news, 1)
	assert.Equal(t, "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena", news[0].Title)
	assert.Equal(t, "The Vatican has issued new rules radically reforming its process for evaluating faith-based supernatural phenomena like visions of the Virgin Mary or stigmata.", news[0].Description)
	assert.Equal(t, "2024-05-17T14:58:52Z", news[0].PubDate)
}
