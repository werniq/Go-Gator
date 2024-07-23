package parsers

import (
	"encoding/xml"
	"gogator/cmd/types"
	"io"
	"net/http"
)

// XMLParser is a struct used to parse article data from XML feeds.
//
// It implements the NewsParser interface, which requires a Parse method.
//
// This struct is used to dynamically parse feeds containing article data.
// The Parse method returns a successfully decoded array of news articles, or an error if parsing fails
type XMLParser struct {
	Source string
}

// Parse fetches and parses XML formatted data from the feed endpoint stored in the XMLParser's Source field.
//
// It retrieves the data from the specified source's endpoint, decodes the XML into a structured format,
// and assigns the source as the publisher for each article.
//
// Returns a slice of parsed news articles and an error, if any
func (xp XMLParser) Parse() ([]types.News, error) {
	res, err := http.Get(sourceToEndpoint[xp.Source])
	if err != nil {
		return nil, err
	}

	var news []types.RSS

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(body, &news)
	if err != nil {
		return nil, err
	}

	articles := news[0].Channel.Items

	for i := 0; i <= len(articles)-1; i++ {
		articles[i].Publisher = xp.Source
	}

	return articles, nil
}
