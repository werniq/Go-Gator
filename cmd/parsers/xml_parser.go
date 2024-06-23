package parsers

import (
	"encoding/xml"
	"io"
	"net/http"
	"newsaggr/cmd/types"
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

	for _, article := range news[0].Channel.Items {
		article.Publisher = xp.Source
	}

	return news[0].Channel.Items, nil
}
