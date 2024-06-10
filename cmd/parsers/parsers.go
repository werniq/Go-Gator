package parsers

import (
	"errors"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"strings"
)

// Parser interface will be used to implement parsers
type Parser interface {
	Parse() []types.News
}

var (
	sourceToFile = map[string]string{
		"washingtontimes": "washington-times.xml",
		"abc":             "abc.xml",
		"bbc":             "bbc.xml",
		"usatoday":        "usa-today.html",
		"nbc":             "nbc-news.json",
	}
	ErrInvalidSource = errors.New("this source is not supported. Supported sources are [abc, bbc, nbc, usatoday, washingtontimes]")
)

// ParseBySource returns all news in particular source. If source is equal to "all", news will be
// retrieved from all sources
func ParseBySource(source string) []types.News {
	var news []types.News

	g := GoGatorParsingFactory{}
	sourceToParser := map[string]Parser{
		"nbc":             g.CreateJsonParser("nbc"),
		"usatoday":        g.CreateHtmlParser("usatoday"),
		"abc":             g.CreateXmlParser("abc"),
		"bbc":             g.CreateXmlParser("bbc"),
		"washingtontimes": g.CreateXmlParser("washingtontimes"),
	}

	if source == "all" {
		for _, p := range sourceToParser {
			news = append(news, p.Parse()...)
		}
	} else {
		splitSources := strings.Split(source, ",")
		for _, source := range splitSources {
			if p, exists := sourceToParser[source]; exists {
				news = append(news, p.Parse()...)
			} else {
				logger.ErrorLogger.Fatalln(ErrInvalidSource)
			}
		}
	}

	return news
}
