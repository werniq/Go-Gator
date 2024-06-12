package parsers

import (
	"newsAggr/cmd/filters"
	"newsAggr/cmd/types"
	"strings"
)

// Parser interface will be used to implement parsers
type Parser interface {
	Parse() ([]types.News, error)
}

var (
	// g and f are Parsing and instruction factories.
	// These are custom types which will be used for parsers initialization for
	// different data formats, and applying filters to retrieved data
	g ParsingFactory
	f filters.InstructionFactory

	// sourceToFile maps source names (as strings) to their corresponding filenames
	sourceToFile = map[string]string{
		"washingtontimes": "washington-times.xml",
		"abc":             "abc.xml",
		"bbc":             "bbc.xml",
		"usatoday":        "usa-today.html",
		"nbc":             "nbc-news.json",
	}

	// sourceToParser maps source names to their corresponding parser instances
	// The parser are created using ParsingFactory
	sourceToParser = map[string]Parser{
		"nbc":             g.CreateJsonParser("nbc"),
		"usatoday":        g.CreateHtmlParser("usatoday"),
		"abc":             g.CreateXmlParser("abc"),
		"bbc":             g.CreateXmlParser("bbc"),
		"washingtontimes": g.CreateXmlParser("washingtontimes"),
	}
)

// ParseBySource returns all news in particular source. If source is equal to "all", news will be
// retrieved from all sources
func ParseBySource(source string) ([]types.News, error) {
	var news []types.News

	if source == "" {
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
			}
		}
	}

	return news, nil
}
