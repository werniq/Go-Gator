package parsers

import (
	"encoding/xml"
	"gogator/cmd/types"
	"io"
	"net/http"
)

type XMLParser struct {
	Source string
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
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
