package parsers

import (
	"io"
	"newsaggr/cmd/types"
	"os"
	"strings"
	"sync"
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

	// sourceToEndpoint maps source names (as strings) to their corresponding filenames
	sourceToEndpoint = map[string]string{
		WashingtonTimes: "https://www.washingtontimes.com/rss/headlines/news/world",
		ABC:             "https://abcnews.go.com/abcnews/internationalheadlines",
		BBC:             "https://feeds.bbci.co.uk/news/rss.xml",
		UsaToday:        "https://usatoday.com",
		//NbcNews:         "nbc-news.json",
	}

	// sourceToParser maps source names to their corresponding parser instances
	// The parser are created using ParsingFactory
	sourceToParser = map[string]Parser{
		//NbcNews:         g.CreateJsonParser("nbc-news.json"),
		UsaToday:        g.CreateHtmlParser(UsaToday),
		ABC:             g.CreateXmlParser(ABC),
		BBC:             g.CreateXmlParser(BBC),
		WashingtonTimes: g.CreateXmlParser(WashingtonTimes),
	}

	// PathToDataDir is used to join path from current working directory to data dir
	PathToDataDir = "\\cmd\\parsers\\data\\"

	// When testing change to:
	//PathToDataDir = "\\data\\"

)

// Define Sources
const (
	NbcNews         = "nbc"
	UsaToday        = "usatoday"
	ABC             = "abc"
	BBC             = "bbc"
	WashingtonTimes = "washingtontimes"

	AllSources = ""

	sourcesFile = "sources.json"

	cwdPath = "C:\\Users\\Oleksandr Matviienko\\GolandProjects\\newsAggr-2\\Go-Gator"
)

// extractFileData reads data from file $filename and returns its content
func extractFileData(filename string) ([]byte, error) {
	file, err := os.Open(cwdPath + PathToDataDir + filename)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ParseBySource returns all news in particular source. If source is equal to "all", news will be
// retrieved from all sources
func ParseBySource(source string) ([]types.News, error) {
	var (
		news       []types.News
		wg         sync.WaitGroup
		mu         sync.Mutex
		errChannel = make(chan error, 1)
	)

	collectNews := func(p Parser) {
		defer wg.Done()
		n, err := p.Parse()
		if err != nil {
			select {
			case errChannel <- err:
			default:
			}
			return
		}
		mu.Lock()
		news = append(news, n...)
		mu.Unlock()
	}

	if source == "" {
		for _, p := range sourceToParser {
			wg.Add(1)
			go collectNews(p)
		}
	} else {
		sources := strings.Split(source, ",")
		for _, s := range sources {
			if p, exists := sourceToParser[s]; exists {
				wg.Add(1)
				go collectNews(p)
			}
		}
	}

	wg.Wait()
	close(errChannel)

	if err, ok := <-errChannel; ok {
		return nil, err
	}

	return news, nil
}
