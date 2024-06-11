package parsers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
	"testing"
)

func TestJsonParser_ParseWithArgs(t *testing.T) {
	params := &types.FilteringParams{
		Keywords: "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena",
	}

	parser := JsonParser{
		"nbc",
	}

	news, err := parser.Parse()
	assert.Nil(t, err)
	news = ApplyParams(news, params)

	assert.Len(t, news, 1)
	assert.Equal(t, "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena", news[0].Title)
	assert.Equal(t, "The Vatican has issued new rules radically reforming its process for evaluating faith-based supernatural phenomena like visions of the Virgin Mary or stigmata.", news[0].Description)
	assert.Equal(t, "2024-05-17T14:58:52Z", news[0].PubDate)
}

func TestJsonParser_Parse(t *testing.T) {
	jsonParser := JsonParser{
		"nbc",
	}

	news, err := jsonParser.Parse()
	assert.Nil(t, err, fmt.Sprintf("Expected: %v, Got: %v", nil, err))

	assert.NotEqual(t, news, nil)
	assert.Equal(t, len(news), 100)
}
