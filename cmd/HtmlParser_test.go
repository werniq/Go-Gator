package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	params := ParsingParams{
		Keywords: "Israeli and Hamas leaders are suspected of committing war crimes during the continuing war in Gaza",
	}
	parser := HtmlParser{}
	news := parser.Parse(params)

	assert.Len(t, news, 1)
	assert.Equal(t, "Israeli and Hamas leaders are suspected of committing war crimes during the continuing war in Gaza", news[0].Description)
	assert.Equal(t, "What ICC arrest warrants mean for Israel and Hamas", news[0].Title)
	assert.Equal(t, "2 hrs ago", news[0].PubDate)
}
