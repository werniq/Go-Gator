package parsers

import (
	"errors"
	"newsAggr/cmd/types"
	"strings"
)

// Parser interface will be used to implement parsers
type Parser interface {
	Parse() ([]types.News, error)
}

var (
	g            ParsingFactory
	sourceToFile = map[string]string{
		"washingtontimes": "washington-times.xml",
		"abc":             "abc.xml",
		"bbc":             "bbc.xml",
		"usatoday":        "usa-today.html",
		"nbc":             "nbc-news.json",
	}
	sourceToParser = map[string]Parser{
		"nbc":             g.CreateJsonParser("nbc"),
		"usatoday":        g.CreateHtmlParser("usatoday"),
		"abc":             g.CreateXmlParser("abc"),
		"bbc":             g.CreateXmlParser("bbc"),
		"washingtontimes": g.CreateXmlParser("washingtontimes"),
	}
	ErrInvalidSource = errors.New("this source is not supported. Supported sources are [abc, bbc, nbc, usatoday, washingtontimes]")
)

// ParseBySource returns all news in particular source. If source is equal to "all", news will be
// retrieved from all sources
func ParseBySource(source string) ([]types.News, error) {
	var news []types.News

	if source == "all" {
		for _, p := range sourceToParser {
			tmp, err := p.Parse()
			if err != nil {
				return nil, err
			}
			news = append(news, tmp...)
		}
	} else {
		splitSources := strings.Split(source, ",")
		for _, source := range splitSources {
			if p, exists := sourceToParser[source]; exists {
				tmp, err := p.Parse()
				if err != nil {
					return nil, err
				}
				news = append(news, tmp...)
			} else {
				return nil, ErrInvalidSource
			}
		}
	}

	return news, nil
}
