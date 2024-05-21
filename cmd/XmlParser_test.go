package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXMLParser_Parse(t *testing.T) {
	params := ParsingParams{
		Keywords: "Ukraine and Russia exchange drone attacks while Russia continues its push in the east",
	}
	parser := XMLParser{}
	news := parser.Parse(params)

	assert.Len(t, news, 1)
	assert.Equal(t, "Ukraine and Russia exchange drone attacks while Russia continues its push in the east", news[0].Title)
	assert.Equal(t, "At least 10 people were reported killed in attacks in Ukraine&rsquo;s war-ravaged northeast on Sunday, as Russia pushed ahead with its renewed offensive", news[0].Description)
	assert.Equal(t, "Sun, 19 May 2024 09:02:27", news[0].PubDate)
}
